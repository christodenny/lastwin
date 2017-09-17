function search(event) {
  var searchBar = document.getElementById("search-bar");
  event.preventDefault();
  window.location.href = "/" + encodeURIComponent(searchBar.value);
}

var searchForm = document.getElementById("search-form");
searchForm.addEventListener("submit", search);
