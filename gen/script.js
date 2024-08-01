window.addEventListener("load", start)


function start () {
    console.log("hello world!")
    hljs.highlightAll();

    const menuButton = document.getElementById("menu")
    const nav = document.getElementById("navbar")
    menuButton.addEventListener("click", () => {
        if (nav.classList.contains("open")) {
            nav.classList.remove("open")
        } else {
            nav.classList.add("open")
        }

    })
}