window.addEventListener("load", start)



function start () {
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

        for (const innerCodeBlock of innerCodeBlocks) {
            codeBlocks.push(innerCodeBlock)
        }

    }

    const codes = []

    for (const codeBlock of codeBlocks) {
        console.log(codeBlock)
        const code = codeBlock.innerText
        console.log(code)
        codes.push(code)
    }

    console.log(codes)

    console.log("highlighting...")
    hljs.highlightAll();

    for (let i = 0; i < codeBlocks.length; i++) {
        const codeBlock = codeBlocks[i]
        const code = codes[i]

        const button = document.createElement("button")
        button.classList.add("copy-button")

        button.addEventListener("click", () => {
            console.log("yo!")

            navigator.clipboard.writeText(code).then(() => {
                alert("Copied code to clipboard.")
            })

        }, false)

        codeBlock.appendChild(button)
    }


}