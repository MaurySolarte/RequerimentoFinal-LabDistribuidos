const icons = document.querySelectorAll('#reacciones [data-reaccion]');

icons.forEach(icon => {
  icon.addEventListener('click', () => {

    // Quitar clase active a todos
    icons.forEach(i => i.classList.remove('active'));

    // Agregar clase active al que se clickeó
    icon.classList.add('active');

    // Solo para mostrar en consola
    console.log(`Reacción seleccionada: ${icon.dataset.reaction}`);
  });
});

function handleReactionClick(icon) {

  // Efecto glow dorado
  icon.classList.add("gold-glow");
  icon.classList.add('active');

  setTimeout(() => {
    icon.classList.remove("gold-glow");
    icon.classList.remove('active');
  }, 300);

}

function mostrarBurbujaReaccion(emoji) {
  const contenedor = document.getElementById('mensajes');
  if (!contenedor) {
    return;
  }

  const emptyState = contenedor.querySelector('.empty-state');
  if (emptyState) {
    emptyState.remove();
  }

  const burbuja = document.createElement('div');
  burbuja.className = 'reaccion-burbuja';
  burbuja.textContent = emoji;
  contenedor.appendChild(burbuja);

  setTimeout(() => {
    burbuja.remove();
    if (contenedor.children.length === 0) {
      const estadoVacio = document.createElement('span');
      estadoVacio.className = 'empty-state';
      estadoVacio.textContent = 'Sin reacciones aún';
      contenedor.appendChild(estadoVacio);
    }
  }, 3250);
}


