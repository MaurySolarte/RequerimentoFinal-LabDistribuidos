    let clienteChat = null 
    
    function setConectado (conectado){
        document.getElementById('btnConectar').disabled = conectado;
        document.getElementById('btnDesconectar').disabled = !conectado;
    }
    
    
    function conectar(){
        const socket = new SockJS('http://localhost:5000/ws');
        clienteChat = Stomp.over(socket);
        
        clienteChat.connect({},onConnected);
    }
    function onConnected(frame){
        console.log("Conectado a: " + frame);
        clienteChat.subscribe('/brokerDeReacciones/reaccionesPorCancion', recibirReaccion);
        setConectado(true);
    }
    
    function recibirReaccion(reaccion){
        console.log("Reaccion reenviada: ", reaccion.body);
        const contenedor = document.getElementById("mensajes");

        //Crear Burbuja 
        const burbuja = document.createElement("div");
        burbuja.classList.add("reaccion-burbuja");

        let icono = "";
        switch(reaccion.body){
            case "like": 
                icono = '<i class="fa-solid fa-thumbs-up icon-like"></i>';
                break;
            case "heart":
                icono = '<i class="fa-solid fa-heart icon-heart"></i>';
                break;

            case "sad":
                icono = '<i class="fa-solid fa-face-sad-tear icon-sad"></i>';
                break;
            
            case "fun":
                icono = '<i class="fa-solid fa-face-laugh icon-fun"></i>';
                break;
                
        }

        console.log("Icono: ", icono);
        console.log("Burbuja: ", burbuja);


        burbuja.innerHTML = icono;
        contenedor.appendChild(burbuja);

        setTimeout(() => {
            burbuja.remove();
        }, 3000);
    }

function enviarReaccionServidor(event){
    const icon = event.target;
    handleReactionClick(icon);
    
    const reaccion = event.target.dataset.reaction;
    if (clienteChat && clienteChat.connected){
        clienteChat.send("/app/enviarReaccion", {}, reaccion);

    }else{
        alert("No estás conectado.");
    }

}

function desconectar(){
    if (clienteChat !== null){
        clienteChat.disconnect(() => {
            setConectado(false);
        });
        clienteChat = null;
    }
}