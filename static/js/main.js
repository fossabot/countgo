document.addEventListener("DOMContentLoaded", function () {
    // ping golang unique visitor counter
    getRequest(location.origin + "/count", function (request) {
        var response = request.currentTarget.response || request.target.responseText;
        var counter = document.getElementById("counter_text");

        if (counter != null) {
            document.getElementById("counter_text")
                .innerHTML = "[" + response + " unique visitors]";
        }
    })
});

function getRequest(url, success) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', url);
    xhr.onload = success;
    xhr.send();
    return xhr;
}

$('#back_to_home1').on('click', function () {
    window.location.href = 'index.html';
});

$('#back_to_home2').on('click', function () {
    window.location.href = 'index.html';
});

$('#yt_text').on('click', function () {
    window.location.href = "music";
});
