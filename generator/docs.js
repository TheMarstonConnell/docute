window.addEventListener("load", init)

function init() {
    console.log("Hello world!")

    let dropDowns = document.getElementsByClassName("dropdown")

    for (const dropDown of dropDowns) {
        dropDown.firstChild.addEventListener("click", () => {
            if (dropDown.classList.contains("open")) {
                dropDown.classList.remove("open")
            } else {
                dropDown.classList.add("open")
            }
        })
    }

    let links = document.getElementsByTagName("a")

    for (const link of links) {
        const page = link.dataset.page
        console.log(page)



        link.addEventListener("click", () => {
            const pageData = document.getElementById("page")
            for (const l of links) {
                    l.classList.remove("active")
            }
            link.classList.add("active")
            fetch(page)
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok ' + response.statusText);
                    }
                    return response.text(); // Get the response body as a string
                })
                .then(data => {
                    console.log('Response body as a string:', data);
                    pageData.innerHTML = data
                    hljs.highlightAll();
                })
                .catch(error => {
                    console.error('There was a problem with the fetch operation:', error);
                });

        })
    }
}