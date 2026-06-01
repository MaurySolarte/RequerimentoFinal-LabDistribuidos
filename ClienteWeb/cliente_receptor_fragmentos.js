const URL_SERVICIO_STREAMING = 'http://localhost:8080';

const reproductorFragmentos = (() => {
  let audioElemento = null;
  let mediaSource = null;
  let fuenteBuffer = null;
  let colaFragmentos = [];
  let procesandoCola = false;
  let controladorAbortado = null;
  let promesaTermino = null;
  let urlObjeto = null;

  function inicializar(elementoAudio) {
    audioElemento = elementoAudio;
  }

  async function reproducir(audioId, referenciaAudio) {
    detener();

    if (!audioElemento) {
      throw new Error('No se inicializo el elemento de audio');
    }

    controladorAbortado = new AbortController();
    mediaSource = new MediaSource();
    colaFragmentos = [];
    procesandoCola = false;

    urlObjeto = URL.createObjectURL(mediaSource);
    audioElemento.src = urlObjeto;

    const listaPromesas = [];
    const promesaFuente = new Promise((resolve, reject) => {
      mediaSource.addEventListener('sourceopen', () => {
        try {
          const mimeType = 'audio/mpeg';
          if (!MediaSource.isTypeSupported(mimeType)) {
            throw new Error('El navegador no soporta MediaSource para audio/mpeg');
          }

          fuenteBuffer = mediaSource.addSourceBuffer(mimeType);
          fuenteBuffer.mode = 'sequence';
          fuenteBuffer.addEventListener('updateend', procesarCola);
          resolve();
        } catch (error) {
          reject(error);
        }
      }, { once: true });

      mediaSource.addEventListener('error', () => reject(new Error('Error inicializando MediaSource')), { once: true });
    });

    listaPromesas.push(promesaFuente);
    promesaTermino = obtenerFragmentos(referenciaAudio, controladorAbortado.signal)
      .then(async () => {
        await esperarColaVacia();
        if (mediaSource && mediaSource.readyState === 'open') {
          try {
            mediaSource.endOfStream();
          } catch (error) {
            console.warn('No fue posible cerrar el MediaSource:', error);
          }
        }
      })
      .catch(error => {
        if (error.name !== 'AbortError') {
          console.error('Error en streaming:', error);
          lanzarEventoError(error);
        }
      });

    await Promise.all(listaPromesas);
    return promesaTermino;
  }

  function detener() {
    if (controladorAbortado) {
      controladorAbortado.abort();
      controladorAbortado = null;
    }

    colaFragmentos = [];
    procesandoCola = false;

    if (fuenteBuffer) {
      try {
        fuenteBuffer.removeEventListener('updateend', procesarCola);
      } catch (error) {
        console.warn('No fue posible retirar el listener del SourceBuffer:', error);
      }
      fuenteBuffer = null;
    }

    if (audioElemento) {
      audioElemento.removeAttribute('src');
      audioElemento.load();
    }

    if (urlObjeto) {
      URL.revokeObjectURL(urlObjeto);
      urlObjeto = null;
    }

    mediaSource = null;
  }

  function pausar() {
    if (audioElemento && !audioElemento.paused) {
      audioElemento.pause();
    }
  }

  async function reanudar() {
    if (audioElemento) {
      await audioElemento.play();
    }
  }

  async function obtenerFragmentos(referenciaAudio, signal) {
    const cuerpoPeticion = construirPeticionGrpc(referenciaAudio);
    const respuesta = await fetch(`${URL_SERVICIO_STREAMING}/servicios.AudioService/enviarCancionMedianteStream`, {
      method: 'POST',
      headers: {
        'content-type': 'application/grpc-web+proto',
        'x-grpc-web': '1',
        'x-user-agent': 'grpc-web-javascript/0.1'
      },
      body: cuerpoPeticion,
      signal
    });

    if (!respuesta.ok && respuesta.status !== 200) {
      const texto = await respuesta.text();
      throw new Error(`El servidor de streaming respondió ${respuesta.status}: ${texto}`);
    }

    if (!respuesta.body) {
      throw new Error('La respuesta de streaming no contiene un cuerpo legible');
    }

    const lector = respuesta.body.getReader();
    let bufferAcumulado = new Uint8Array(0);

    while (true) {
      const { value, done } = await lector.read();
      if (done) {
        break;
      }

      bufferAcumulado = concatenarBytes(bufferAcumulado, value);
      const resultado = extraerMensajesGrpcWeb(bufferAcumulado);
      bufferAcumulado = resultado.restante;

      for (const mensaje of resultado.mensajes) {
        const fragmento = decodificarFragmentoCancion(mensaje);
        if (fragmento && fragmento.byteLength > 0) {
          encolarFragmento(fragmento);
        }
      }
    }
  }

  function encolarFragmento(fragmento) {
    colaFragmentos.push(fragmento);
    procesarCola();
  }

  function procesarCola() {
    if (!fuenteBuffer || procesandoCola || fuenteBuffer.updating || colaFragmentos.length === 0) {
      return;
    }

    procesandoCola = true;
    const siguiente = colaFragmentos.shift();

    try {
      fuenteBuffer.appendBuffer(siguiente);
    } catch (error) {
      console.error('No se pudo anexar el fragmento de audio:', error);
    } finally {
      procesandoCola = false;
    }
  }

  async function esperarColaVacia() {
    while (colaFragmentos.length > 0 || (fuenteBuffer && fuenteBuffer.updating)) {
      await esperar(50);
    }
  }

  function construirPeticionGrpc(referenciaAudio) {
    const texto = String(referenciaAudio ?? '');
    const bytesTitulo = codificarCadena(1, texto);
    const bytesFormato = codificarCadena(2, '');
    const mensaje = concatenarBytes(bytesTitulo, bytesFormato);
    return empaquetarGrpcWeb(mensaje);
  }

  function empaquetarGrpcWeb(mensaje) {
    const envoltura = new Uint8Array(5 + mensaje.length);
    envoltura[0] = 0;
    const longitud = mensaje.length;
    envoltura[1] = (longitud >>> 24) & 0xff;
    envoltura[2] = (longitud >>> 16) & 0xff;
    envoltura[3] = (longitud >>> 8) & 0xff;
    envoltura[4] = longitud & 0xff;
    envoltura.set(mensaje, 5);
    return envoltura;
  }

  function codificarCadena(numeroCampo, valor) {
    if (!valor) {
      return new Uint8Array(0);
    }

    const texto = new TextEncoder().encode(valor);
    const salida = new Uint8Array(1 + longitudVarint(texto.length) + texto.length);
    salida[0] = (numeroCampo << 3) | 2;
    const longitudBytes = codificarVarint(texto.length);
    salida.set(longitudBytes, 1);
    salida.set(texto, 1 + longitudBytes.length);
    return salida;
  }

  function decodificarFragmentoCancion(mensaje) {
    if (!mensaje || mensaje.length === 0) {
      return new Uint8Array(0);
    }

    let indice = 0;
    while (indice < mensaje.length) {
      const etiqueta = mensaje[indice++];
      const numeroCampo = etiqueta >>> 3;
      const tipo = etiqueta & 0x07;
      if (numeroCampo === 1 && tipo === 2) {
        const longitud = leerVarint(mensaje, indice);
        indice += longitud.bytesLeidos;
        const inicio = indice;
        const fin = inicio + longitud.valor;
        return mensaje.slice(inicio, fin);
      }

      if (tipo === 2) {
        const longitud = leerVarint(mensaje, indice);
        indice += longitud.bytesLeidos + longitud.valor;
      } else if (tipo === 0) {
        const longitud = leerVarint(mensaje, indice);
        indice += longitud.bytesLeidos;
      } else if (tipo === 5) {
        indice += 4;
      } else if (tipo === 1) {
        indice += 8;
      } else {
        break;
      }
    }

    return new Uint8Array(0);
  }

  function extraerMensajesGrpcWeb(buffer) {
    const mensajes = [];
    let indice = 0;

    while (indice + 5 <= buffer.length) {
      const tipo = buffer[indice];
      const longitud = (buffer[indice + 1] << 24) | (buffer[indice + 2] << 16) | (buffer[indice + 3] << 8) | buffer[indice + 4];
      if (indice + 5 + longitud > buffer.length) {
        break;
      }

      if ((tipo & 0x80) === 0) {
        mensajes.push(buffer.slice(indice + 5, indice + 5 + longitud));
      }
      indice += 5 + longitud;
    }

    return {
      mensajes,
      restante: buffer.slice(indice)
    };
  }

  function leerVarint(bytes, indice) {
    let valor = 0;
    let desplazamiento = 0;
    let bytesLeidos = 0;
    let actual = 0;

    do {
      actual = bytes[indice + bytesLeidos];
      valor |= (actual & 0x7f) << desplazamiento;
      desplazamiento += 7;
      bytesLeidos += 1;
    } while (actual >= 0x80 && indice + bytesLeidos < bytes.length);

    return { valor, bytesLeidos };
  }

  function codificarVarint(valor) {
    const bytes = [];
    let numero = valor >>> 0;
    while (numero >= 0x80) {
      bytes.push((numero & 0x7f) | 0x80);
      numero >>>= 7;
    }
    bytes.push(numero);
    return new Uint8Array(bytes);
  }

  function longitudVarint(valor) {
    let longitud = 1;
    let numero = valor >>> 0;
    while (numero >= 0x80) {
      numero >>>= 7;
      longitud += 1;
    }
    return longitud;
  }

  function concatenarBytes(a, b) {
    const resultado = new Uint8Array(a.length + b.length);
    resultado.set(a, 0);
    resultado.set(b, a.length);
    return resultado;
  }

  function esperar(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  function lanzarEventoError(error) {
    window.dispatchEvent(new CustomEvent('error-streaming-audio', { detail: error }));
  }

  return {
    inicializar,
    reproducir,
    detener,
    pausar,
    reanudar
  };
})();

window.reproductorFragmentos = reproductorFragmentos;
