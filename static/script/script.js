document.addEventListener("htmx:afterRequest", function(evt) {
    // 1. Chiudi la modale DaisyUI
    console.log(evt)
    console.log(evt.detail.elt.id)
    if (evt.detail.elt.id === "creationForm" && evt.detail.successful === true) {

        document.getElementById('my_modal_1').close();
        const toast = document.getElementById('toast');
        // window.location.reload();
        toast.classList.remove('hidden'); // Rimuovi la classe 'hidden' per mostrare il toast

        // // 3. Mantieni il toast visibile per 2 secondi e poi ricarica la pagina
        setTimeout(() => {
            toast.classList.add('hidden'); 
        }, 3000);
    }
});