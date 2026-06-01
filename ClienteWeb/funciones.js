let clienteChat = null;
let reproduciendo = false;
let pausado = false;
let suscripcionesCancion = [];

const cancionActual = {
    title: document.getElementById('song-title')?.textContent?.trim() || 'Cancion actual',
    artist: document.getElementById('song-artist')?.textContent?.trim() || 'Artista actual',
    key: 'general'
};

function setConectado(conectado) {
    document.getElementById('btnConectar').disabled = conectado;
    document.getElementById('btnDesconectar').disabled = !conectado;
}

function conectar() {
    const socket = new SockJS('http://localhost:5000/ws');
    clienteChat = Stomp.over(socket);
    clienteChat.debug = null;

    clienteChat.connect({}, onConnected);
}

function onConnected(frame) {
    console.log('Conectado a: ' + frame);

    setConectado(true);
}

function getSongKey() {
    return cancionActual.key || 'general';
}

function getNickname() {
    const nickname = document.getElementById('nickname').value.trim();
    return nickname || 'Anónimo';
}

function renderMovimiento(listaId, nickname, etiqueta, claseExtra) {
    const contenedor = document.getElementById(listaId);
    if (!contenedor) {
        return;
    }

    const emptyState = contenedor.querySelector('.empty-state');
    if (emptyState) {
        emptyState.remove();
    }

    const tarjeta = document.createElement('div');
    tarjeta.className = `nick-item${claseExtra ? ` ${claseExtra}` : ''}`;
    tarjeta.innerHTML = `
        <span class="dot"></span>
        <span>${escapeHtml(nickname)}</span>
        <span class="ts">${nowTime()}</span>
        <span class="label">${escapeHtml(etiqueta)}</span>
    `;

    contenedor.prepend(tarjeta);

    while (contenedor.children.length > 6) {
        contenedor.removeChild(contenedor.lastElementChild);
    }
}

function recibirInicioReproduccion(reaccion) {
    const payload = parsePayload(reaccion.body);
    const nickname = payload.nickname || 'Anónimo';
    const etiqueta = payload.action === 'resume' ? 'Reanudó la reproducción' : 'Comenzó a reproducir';
    renderMovimiento('playing-list', nickname, etiqueta, '');
    // Solo cambiar el estado local si el mensaje corresponde a este cliente
    if (nickname === getNickname()) {
        reproduciendo = true;
        pausado = false;
        const icon = document.getElementById('play-icon');
        if (icon) icon.className = 'fa-solid fa-pause';
    }
}

function recibirFinReproduccion(reaccion) {
    const payload = parsePayload(reaccion.body);
    const nickname = payload.nickname || 'Anónimo';
    const etiqueta = payload.action === 'finish' ? 'Finalizó la reproducción' : 'Pausó la reproducción';
    renderMovimiento('stopped-list', nickname, etiqueta, 'stop');
    // Solo cambiar el estado local si el mensaje corresponde a este cliente
    if (nickname === getNickname()) {
        reproduciendo = false;
        pausado = true;
        const icon = document.getElementById('play-icon');
        if (icon) icon.className = 'fa-solid fa-play';
    }
}

function recibirReaccion(reaccion) {
    const payload = parsePayload(reaccion.body);
    const reaccionTipo = payload.reaction || payload.action || reaccion.body;
    const contenedor = document.getElementById('mensajes');

    const emptyState = contenedor.querySelector('.empty-state');
    if (emptyState) {
        emptyState.remove();
    }

    const burbuja = document.createElement('div');
    burbuja.classList.add('reaccion-burbuja');

    let icono = '';
    switch (reaccionTipo) {
        case 'like':
            icono = '<i class="fa-solid fa-thumbs-up icon-like"></i>';
            break;
        case 'heart':
            icono = '<i class="fa-solid fa-heart icon-heart"></i>';
            break;
        case 'sad':
            icono = '<i class="fa-solid fa-face-sad-tear icon-sad"></i>';
            break;
        case 'fun':
            icono = '<i class="fa-solid fa-face-laugh icon-fun"></i>';
            break;
        default:
            icono = '<i class="fa-solid fa-star"></i>';
            break;
    }

    burbuja.innerHTML = icono;
    contenedor.appendChild(burbuja);

    setTimeout(() => {
        burbuja.remove();
    }, 3000);
}

function enviarReaccionServidor(event) {
    const icon = event.currentTarget;
    handleReactionClick(icon);

    const reaccion = icon.dataset.reaction;
    if (clienteChat && clienteChat.connected && reproduciendo) {
        clienteChat.send('/app/enviarReaccion', { 'content-type': 'application/json' }, JSON.stringify({
            nickname: getNickname(),
            songKey: getSongKey(),
            songTitle: cancionActual.title,
            songArtist: cancionActual.artist,
            reaction: reaccion
        }));
    } else {
        alert('Debes estar reproduciendo la canción para enviar reacciones.');
    }
}

function togglePlay() {
    const icon = document.getElementById('play-icon');

    if (!clienteChat || !clienteChat.connected) {
        alert('No estás conectado.');
        return;
    }

    if (!reproduciendo) {
        const accion = pausado ? 'resume' : 'start';
        suscribirseACanalCancion();
        clienteChat.send('/app/iniciarReproduccion', { 'content-type': 'application/json' }, JSON.stringify({
            nickname: getNickname(),
            songKey: getSongKey(),
            songTitle: cancionActual.title,
            songArtist: cancionActual.artist,
            action: accion
        }));
        icon.className = 'fa-solid fa-pause';
        reproduciendo = true;
        pausado = false;
        return;
    }

    clienteChat.send('/app/detenerReproduccion', { 'content-type': 'application/json' }, JSON.stringify({
        nickname: getNickname(),
        songKey: getSongKey(),
        songTitle: cancionActual.title,
        songArtist: cancionActual.artist,
        action: 'pause'
    }));
    icon.className = 'fa-solid fa-play';
    reproduciendo = false;
    pausado = true;
    desuscribirseDeCanalCancion();
}

function suscribirseACanalCancion() {
    if (!clienteChat || !clienteChat.connected) {
        return;
    }

    desuscribirseDeCanalCancion();

    const songKey = getSongKey();
    suscripcionesCancion.push(clienteChat.subscribe(`/brokerDeReacciones/canciones/${songKey}/reacciones`, recibirReaccion));
    suscripcionesCancion.push(clienteChat.subscribe(`/brokerDeReacciones/canciones/${songKey}/reproducciones/inicio`, recibirInicioReproduccion));
    suscripcionesCancion.push(clienteChat.subscribe(`/brokerDeReacciones/canciones/${songKey}/reproducciones/fin`, recibirFinReproduccion));
}

function desuscribirseDeCanalCancion() {
    suscripcionesCancion.forEach(subscription => {
        if (subscription && typeof subscription.unsubscribe === 'function') {
            subscription.unsubscribe();
        }
    });

    suscripcionesCancion = [];
}

function parsePayload(body) {
    try {
        return JSON.parse(body);
    } catch (error) {
        return { reaction: body, nickname: null, action: null };
    }
}

function escapeHtml(value) {
    return String(value)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

function nowTime() {
    const date = new Date();
    return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
}

function desconectar() {
    if (clienteChat !== null) {
        desuscribirseDeCanalCancion();
        clienteChat.disconnect(() => {
            setConectado(false);
        });
        clienteChat = null;
    }
}