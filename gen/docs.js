window.addEventListener("load", init)

function updatePage(pageData, link) {
    console.log(link)

    fetch(link)
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

}

function fixA(pageData, linkLocation) {
    const hrefURL = new URL(linkLocation);
    const pageURL = new URL(window.location);
    console.log(hrefURL.host, pageURL.host)
    if (hrefURL.host === pageURL.host) {
        updatePage(pageData, linkLocation)
        return
    }

    console.log("Skipping " + linkLocation)
    window.open(linkLocation, '_blank').focus()
}

function fixAs(page) {


    document.onclick = function (e) {
        e = e ||  window.event;
        const element = e.target || e.srcElement;

        if (element.tagName === 'A' || element.tagName === 'a') {
            fixA(page, element.href);
            return false; // prevent default action and stop event propagation
        }
    };
}

function init() {
    console.log("Hello docute!")
    const page = document.getElementById("page")

    fixAs(page)

    updatePage(page, "README.html")

}