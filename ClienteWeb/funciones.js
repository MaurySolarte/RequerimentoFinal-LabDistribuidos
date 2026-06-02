const URL_SERVICIO_AUDIOS = 'http://localhost:8081';
const URL_SERVICIO_CALLBACK = 'http://localhost:5000/ws';

let clienteChat = null;
let nicknameUsuario = '';
let audioSeleccionado = null;
let audioEnReproduccion = false;
let audioEnPausa = false;
let audioYaInicio = false;
let estadoSilencioso = false;
let tipoSeleccionadoActual = null;
let suscripcionesCanalActual = [];
let tiposAudioDisponibles = [];
let audiosDisponibles = [];

const etiquetasMetadatos = {
  musica: {
    artistaPrincipal: 'Artista Principal',
    album: 'Álbum',
    generoMusical: 'Género Musical'
  },
  podcast: {
    nombrePodcast: 'Nombre del Podcast',
    tituloEpisodio: 'Título del Episodio',
    anfitrion: 'Anfitrión'
  },
  audiolibro: {
    tituloLibro: 'Título del Libro',
    autor: 'Autor',
    narrador: 'Narrador'
  },
  'ruido blanco': {
    tipoSonido: 'Tipo de Sonido',
    fuenteAudio: 'Fuente del Audio',
    usoSugerido: 'Uso Sugerido'
  }
};

document.addEventListener('DOMContentLoaded', iniciarInterfaz);

function iniciarInterfaz() {
  const audioElemento = document.getElementById('audio-reproductor');
  if (window.reproductorFragmentos && audioElemento) {
    window.reproductorFragmentos.inicializar(audioElemento);
  }

  audioElemento?.addEventListener('play', manejarEventoPlay);
  audioElemento?.addEventListener('pause', manejarEventoPause);

  document.getElementById('btnConectar')?.addEventListener('click', conectar);
  document.getElementById('btnDesconectar')?.addEventListener('click', desconectar);
  document.getElementById('btnConfirmarNickname')?.addEventListener('click', confirmarNickname);
  document.getElementById('nicknameModalInput')?.addEventListener('keydown', evento => {
    if (evento.key === 'Enter') {
      confirmarNickname();
    }
  });

  document.getElementById('play-btn')?.addEventListener('click', togglePlay);
  document.querySelectorAll('[data-reaccion]').forEach(boton => {
    boton.addEventListener('click', enviarReaccionServidor);
  });

  ocultarPanelReacciones();
  abrirModalNickname();
  cargarTiposAudio();
}

function abrirModalNickname() {
  const modal = document.getElementById('modal-nickname');
  if (modal) {
    modal.classList.add('visible');
  }
  document.getElementById('nicknameModalInput')?.focus();
}

function cerrarModalNickname() {
  const modal = document.getElementById('modal-nickname');
  if (modal) {
    modal.classList.remove('visible');
  }
}

function confirmarNickname() {
  const entradaModal = document.getElementById('nicknameModalInput');
  const nickname = entradaModal?.value.trim() || document.getElementById('nickname')?.value.trim();

  if (!nickname) {
    alert('Debes escribir un nickname para continuar.');
    return;
  }

  nicknameUsuario = nickname;
  const inputTopBar = document.getElementById('nickname');
  if (inputTopBar) {
    inputTopBar.value = nickname;
  }

  cerrarModalNickname();
  conectar();
}

function setConectado(conectado) {
  const botonConectar = document.getElementById('btnConectar');
  const botonDesconectar = document.getElementById('btnDesconectar');
  const estadoConexion = document.getElementById('estado-conexion');

  if (botonConectar) {
    botonConectar.disabled = conectado;
  }
  if (botonDesconectar) {
    botonDesconectar.disabled = !conectado;
  }
  if (estadoConexion) {
    estadoConexion.textContent = conectado ? 'Conectado' : 'Desconectado';
  }
}

function conectar() {
  if (clienteChat && clienteChat.connected) {
    return;
  }

  if (!nicknameUsuario) {
    nicknameUsuario = document.getElementById('nickname')?.value.trim() || '';
  }

  if (!nicknameUsuario) {
    abrirModalNickname();
    return;
  }

  const socket = new SockJS(URL_SERVICIO_CALLBACK);
  clienteChat = Stomp.over(socket);
  clienteChat.debug = null;
  clienteChat.connect({}, manejarConexion);
}

function manejarConexion(frame) {
  console.log('Conectado a callback:', frame);
  setConectado(true);
  if (audioSeleccionado) {
    suscribirseCanalesAudio(audioSeleccionado.id);
  }
}

function desconectar() {
  desuscribirseCanalesAudio();

  if (clienteChat) {
    clienteChat.disconnect(() => {
      setConectado(false);
    });
    clienteChat = null;
  }
}

function cargarTiposAudio() {
  fetch(`${URL_SERVICIO_AUDIOS}/tipos-audio`)
    .then(respuesta => respuesta.json())
    .then(tipos => {
      tiposAudioDisponibles = Array.isArray(tipos) ? tipos : [];
      renderizarTiposAudio(tiposAudioDisponibles);
      if (tiposAudioDisponibles.length > 0) {
        seleccionarTipoAudio(tiposAudioDisponibles[0]);
      }
    })
    .catch(error => {
      console.error('No se pudieron cargar los tipos de audio:', error);
      mostrarMensajeVacio('lista-tipos', 'No fue posible cargar los tipos de audio');
    });
}

function renderizarTiposAudio(tipos) {
  const contenedor = document.getElementById('lista-tipos');
  if (!contenedor) {
    return;
  }

  contenedor.innerHTML = '';
  tipos.forEach(tipo => {
    const boton = document.createElement('button');
    boton.className = `chip-tipo${tipoSeleccionadoActual === tipo.id ? ' activo' : ''}`;
    boton.textContent = tipo.nombre;
    boton.addEventListener('click', () => seleccionarTipoAudio(tipo));
    contenedor.appendChild(boton);
  });
}

function seleccionarTipoAudio(tipo) {
  tipoSeleccionadoActual = tipo.id;
  renderizarTiposAudio(tiposAudioDisponibles);
  cargarAudiosPorTipo(tipo.id);
}

function cargarAudiosPorTipo(tipoId) {
  fetch(`${URL_SERVICIO_AUDIOS}/audios?tipo_id=${encodeURIComponent(tipoId)}`)
    .then(respuesta => respuesta.json())
    .then(audios => {
      audiosDisponibles = Array.isArray(audios) ? audios : [];
      renderizarAudios(audiosDisponibles);
    })
    .catch(error => {
      console.error('No se pudieron cargar los audios:', error);
      mostrarMensajeVacio('lista-audios', 'No fue posible cargar los audios de este tipo');
    });
}

function renderizarAudios(audios) {
  const contenedor = document.getElementById('lista-audios');
  if (!contenedor) {
    return;
  }

  contenedor.innerHTML = '';
  if (audios.length === 0) {
    mostrarMensajeVacio('lista-audios', 'No hay audios registrados para este tipo');
    return;
  }

  audios.forEach(audio => {
    const boton = document.createElement('button');
    boton.className = `audio-item${audioSeleccionado && audioSeleccionado.id === audio.id ? ' activo' : ''}`;
    boton.innerHTML = `
      <span class="audio-titulo">${escapeHtml(audio.titulo)}</span>
      <span class="audio-subtitulo">#${audio.id} · ${escapeHtml(audio.tipoNombre || '')}</span>
    `;
    boton.addEventListener('click', () => seleccionarAudio(audio));
    contenedor.appendChild(boton);
  });
}

function seleccionarAudio(audio) {
  audioSeleccionado = audio;
  audioYaInicio = false;
  audioEnReproduccion = false;
  audioEnPausa = false;
  actualizarBotonesAudio();
  renderizarAudios(audiosDisponibles);
  cargarDetallesAudio(audio.id);
  suscribirseCanalesAudio(audio.id);
  estadoSilencioso = true;
  if (window.reproductorFragmentos) {
    window.reproductorFragmentos.detener();
  }
  estadoSilencioso = false;
  actualizarEstadoReproduccion('detenido');
}

function cargarDetallesAudio(audioId) {
  fetch(`${URL_SERVICIO_AUDIOS}/audios/${encodeURIComponent(audioId)}`)
    .then(respuesta => respuesta.json())
    .then(audio => {
      mostrarDetallesAudio(audio);
    })
    .catch(error => {
      console.error('No se pudieron cargar los detalles del audio:', error);
      mostrarMensajeVacio('detalle-audio', 'No fue posible cargar los metadatos');
    });
}

function mostrarDetallesAudio(audio) {
  const contenedor = document.getElementById('detalle-audio');
  if (!contenedor) {
    return;
  }

  const campos = construirCamposDetalle(audio);
  const tituloVisible = audio.titulo || audio.title || '';
  const tipoVisible = audio.tipoNombre || audio.typeName || '';

  const tituloElemento = document.getElementById('song-title');
  const artistaElemento = document.getElementById('song-artist');
  if (tituloElemento) {
    tituloElemento.textContent = tituloVisible || 'Audio seleccionado';
  }
  if (artistaElemento) {
    artistaElemento.textContent = tipoVisible || 'Metadatos dinámicos';
  }

  contenedor.innerHTML = `
    <div class="detalle-cabecera">
      <span class="detalle-tipo">${escapeHtml(tipoVisible)}</span>
      <h3>${escapeHtml(tituloVisible)}</h3>
      <p>Archivo: ${escapeHtml(audio.rutaArchivo || audio.filename || '')}</p>
    </div>
    <div class="detalle-matriz">
      ${campos.map(campo => `
        <div class="detalle-campo">
          <span>${escapeHtml(campo.etiqueta)}</span>
          <strong>${escapeHtml(campo.valor)}</strong>
        </div>
      `).join('')}
    </div>
  `;
}

function construirCamposDetalle(audio) {
  const tipo = normalizarTexto(audio.tipoNombre || '');
  const metadatos = audio.metadatos || {};
  const mapaEtiquetas = etiquetasMetadatos[tipo] || {};

  return Object.entries(mapaEtiquetas).map(([clave, etiqueta]) => ({
    etiqueta,
    valor: metadatos[clave] || metadatos[normalizarTexto(etiqueta)] || ''
  }));
}

function togglePlay() {
  if (!audioSeleccionado) {
    alert('Selecciona un audio antes de reproducirlo.');
    return;
  }

  if (!clienteChat || !clienteChat.connected) {
    conectar();
    return;
  }

  if (!audioEnReproduccion) {
    iniciarReproduccionReal();
    return;
  }

  const audioElemento = document.getElementById('audio-reproductor');
  audioElemento?.pause();
}

async function iniciarReproduccionReal() {
  if (!audioSeleccionado) {
    return;
  }

  if (window.reproductorFragmentos) {
    if (!audioYaInicio) {
      estadoSilencioso = true;
      await window.reproductorFragmentos.reproducir(audioSeleccionado.id, audioSeleccionado.id);
      estadoSilencioso = false;
    }
  }

  const audioElemento = document.getElementById('audio-reproductor');
  if (audioElemento) {
    await audioElemento.play();
  }
}

function manejarEventoPlay() {
  if (estadoSilencioso) {
    return;
  }

  if (!audioSeleccionado || !clienteChat || !clienteChat.connected) {
    return;
  }

  if (!audioYaInicio) {
    enviarCallbackReproduccion('nuevo_oyente');
    audioYaInicio = true;
  } else if (audioEnPausa) {
    enviarCallbackReproduccion('reanuda');
  }

  audioEnReproduccion = true;
  audioEnPausa = false;
  actualizarEstadoReproduccion('reproduciendo');
}

function manejarEventoPause() {
  if (estadoSilencioso) {
    return;
  }

  if (!audioSeleccionado || !clienteChat || !clienteChat.connected) {
    return;
  }

  if (audioEnReproduccion) {
    enviarCallbackReproduccion('pausa');
  }

  audioEnReproduccion = false;
  audioEnPausa = true;
  actualizarEstadoReproduccion('pausado');
}

function enviarCallbackReproduccion(tipo) {
  const mensaje = {
    tipo,
    nickname: nicknameUsuario,
    audioId: String(audioSeleccionado?.id || ''),
    titulo: audioSeleccionado?.titulo || '',
    tipoAudio: audioSeleccionado?.tipoNombre || ''
  };

  if (tipo === 'pausa') {
    clienteChat.send('/app/detenerReproduccion', { 'content-type': 'application/json' }, JSON.stringify(mensaje));
    return;
  }

  clienteChat.send('/app/iniciarReproduccion', { 'content-type': 'application/json' }, JSON.stringify(mensaje));
}

function enviarReaccionServidor(evento) {
  const boton = evento.currentTarget;
  handleReactionClick(boton);

  if (!audioSeleccionado || (!audioEnReproduccion && !audioEnPausa)) {
    alert('Solo puedes reaccionar mientras hay un audio en reproducción o en pausa.');
    return;
  }

  if (!clienteChat || !clienteChat.connected) {
    alert('Debes estar conectado para enviar reacciones.');
    return;
  }

  const emoji = boton.dataset.reaccion;
  clienteChat.send('/app/enviarReaccion', { 'content-type': 'application/json' }, JSON.stringify({
    tipo: 'reaccion',
    nickname: nicknameUsuario,
    audioId: String(audioSeleccionado.id),
    emoji,
    reaction: emoji,
    titulo: audioSeleccionado.titulo || ''
  }));
}

function suscribirseCanalesAudio(audioId) {
  if (!clienteChat || !clienteChat.connected) {
    return;
  }

  desuscribirseCanalesAudio();

  suscripcionesCanalActual.push(clienteChat.subscribe(`/brokerDeReacciones/canciones/${audioId}/oyentes`, recibirOyente));
  suscripcionesCanalActual.push(clienteChat.subscribe(`/brokerDeReacciones/canciones/${audioId}/estados`, recibirEstadoReproduccion));
  suscripcionesCanalActual.push(clienteChat.subscribe(`/brokerDeReacciones/canciones/${audioId}/reacciones`, recibirReaccion));
}

function desuscribirseCanalesAudio() {
  suscripcionesCanalActual.forEach(suscripcion => {
    if (suscripcion && typeof suscripcion.unsubscribe === 'function') {
      suscripcion.unsubscribe();
    }
  });
  suscripcionesCanalActual = [];
}

function recibirOyente(mensaje) {
  const payload = parsearPayload(mensaje.body);
  if (payload.tipo === 'pausa' || payload.estado === 'pausó') {
    agregarEntradaActividad('actividad-oyentes', payload.nickname || 'Anónimo', 'Pausó la reproducción', 'entrada');
    return;
  }

  agregarEntradaActividad('actividad-estados', payload.nickname || 'Anónimo', 'Comenzó a reproducir', 'entrada');
}

function recibirEstadoReproduccion(mensaje) {
  const payload = parsearPayload(mensaje.body);
  const estado = payload.estado || payload.tipo || 'estado';
  if (estado === 'reanudó' || estado === 'reanuda') {
    agregarEntradaActividad('actividad-estados', payload.nickname || 'Anónimo', 'Reanudó la reproducción', 'estado');
    return;
  }

  if (estado === 'pausó' || estado === 'pausa') {
    agregarEntradaActividad('actividad-oyentes', payload.nickname || 'Anónimo', 'Pausó la reproducción', 'estado');
  }
}

function recibirReaccion(mensaje) {
  const payload = parsearPayload(mensaje.body);
  mostrarBurbujaReaccion(payload.emoji || payload.reaction || mensaje.body);
}

function agregarEntradaActividad(contenedorId, nickname, texto, claseExtra) {
  const contenedor = document.getElementById(contenedorId);
  if (!contenedor) {
    return;
  }

  eliminarEstadoVacio(contenedor);
  const entrada = document.createElement('div');
  entrada.className = `entrada-actividad ${claseExtra || ''}`.trim();
  entrada.innerHTML = `
    <span class="entrada-punto"></span>
    <strong>${escapeHtml(nickname)}</strong>
    <span>${escapeHtml(texto)}</span>
    <time>${horaActual()}</time>
  `;
  contenedor.prepend(entrada);
  limitarHijos(contenedor, 6);
  mostrarBurbujaReaccion(emoji);
}

function mostrarPanelReacciones() {
  document.getElementById('panel-reacciones')?.classList.remove('oculto');
}

function ocultarPanelReacciones() {
  document.getElementById('panel-reacciones')?.classList.add('oculto');
}

function actualizarBotonesAudio() {
  document.querySelectorAll('#reacciones button').forEach(boton => {
    boton.disabled = !(audioEnReproduccion || audioEnPausa);
  });
}

  mostrarBurbujaReaccion(emoji);
function actualizarEstadoReproduccion(estado) {
  const icono = document.getElementById('play-icon');
  const etiqueta = document.getElementById('estado-reproduccion');

  if (estado === 'reproduciendo') {
    audioEnReproduccion = true;
    audioEnPausa = false;
    if (icono) icono.className = 'fa-solid fa-pause';
    if (etiqueta) etiqueta.textContent = 'Reproduciendo';
  } else if (estado === 'pausado') {
    audioEnReproduccion = false;
    audioEnPausa = true;
    if (icono) icono.className = 'fa-solid fa-play';
    if (etiqueta) etiqueta.textContent = 'Pausado';
  } else {
    audioEnReproduccion = false;
    audioEnPausa = false;
    if (icono) icono.className = 'fa-solid fa-play';
    if (etiqueta) etiqueta.textContent = 'Detenido';
    ocultarPanelReacciones();
  }

  if (estado === 'reproduciendo' || estado === 'pausado') {
    mostrarPanelReacciones();
  }

  actualizarBotonesAudio();
}

function parsearPayload(body) {
  try {
    return JSON.parse(body);
  } catch (error) {
    return { nickname: 'Anónimo', emoji: body, reaction: body };
  }
}

function mostrarMensajeVacio(contenedorId, mensaje) {
  const contenedor = document.getElementById(contenedorId);
  if (!contenedor) {
    return;
  }

  contenedor.innerHTML = `<span class="empty-state">${escapeHtml(mensaje)}</span>`;
}

function eliminarEstadoVacio(contenedor) {
  contenedor.querySelectorAll('.empty-state').forEach(elemento => elemento.remove());
}

function limitarHijos(contenedor, maximo) {
  while (contenedor.children.length > maximo) {
    contenedor.removeChild(contenedor.lastElementChild);
  }
}

function horaActual() {
  const fecha = new Date();
  return `${String(fecha.getHours()).padStart(2, '0')}:${String(fecha.getMinutes()).padStart(2, '0')}`;
}

function escapeHtml(valor) {
  return String(valor)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

function normalizarTexto(valor) {
  return String(valor)
    .trim()
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
}

window.conectar = conectar;
window.desconectar = desconectar;
window.togglePlay = togglePlay;
window.enviarReaccionServidor = enviarReaccionServidor;
