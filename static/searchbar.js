var searchBar = document.getElementById("search-bar");
function search(event) {
  event.preventDefault();
  var dt = new Date();
  var tz = dt.getTimezoneOffset();
  window.location.href = "/" + encodeURIComponent(searchBar.value) + "?tz=" + tz;
}

var req;
function oninput(event) {
  req = new XMLHttpRequest();
  if (!req) {
    return false;
  }
  req.onreadystatechange = autocomplete;
  req.open("GET", "/autocomplete?text=" + encodeURIComponent(searchBar.value), true);
  req.send();
}

var autocompleteBoxes = document.getElementsByClassName("autocomplete-box");
for (var i = 0; i < 3; i++) {
  var element = autocompleteBoxes[i];
  element.addEventListener("click", function(e) {
    return function() {
      var dt = new Date();
      var tz = dt.getTimezoneOffset();
      window.location = "/" + e.textContent + "?tz=" + tz;
    }
  }(element));
}

function autocomplete() {
  if (req.readyState === XMLHttpRequest.DONE && req.status == 200) {
    var data = JSON.parse(req.response);
    for (var i = 0; i < 3; i++) {
      var box = autocompleteBoxes[i];
      if (i < data.length) {
        box.textContent = data[i];
        box.style.visibility = "visible";
      } else {
        box.style.visibility = "hidden";
      }
    }
  }
}

searchBar.addEventListener('input', oninput);

var searchForm = document.getElementById("search-form");
searchForm.addEventListener("submit", search);
