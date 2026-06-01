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


