window.addEventListener("load", start)


function start() {
    console.log("hello world!")

    const menuButton = document.getElementById("menu")
    const nav = document.getElementById("navbar")
    menuButton.addEventListener("click", () => {
        if (nav.classList.contains("open")) {
            nav.classList.remove("open")
        } else {
            nav.classList.add("open")
        }

    })

    const codeBlocks = []

    const allPres = document.getElementsByTagName("pre")

    for (const allPre of allPres) {
        const innerCodeBlocks = allPre.getElementsByTagName("code")

        const code = innerCodeBlocks[0].innerText


        const button = document.createElement("button")
        button.classList.add("copy-button")
        button.ariaLabel = "Copy Code Button"

        button.addEventListener("click", () => {
            console.log("yo!")

            navigator.clipboard.writeText(code).then(() => {
                let x = document.getElementById("snackbar");

                // Add the "show" class to DIV
                x.innerText = "Code copied to clipboard."
                x.className = "show";

                // After 3 seconds, remove the show class from DIV
                setTimeout(function () {
                    x.className = x.className.replace("show", "");
                }, 3000);

            })

        }, false)

        allPre.appendChild(button)

    }

    console.log("highlighting...")
    hljs.highlightAll();




}