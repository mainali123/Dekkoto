let preloader = document.querySelector(".preloader");
let toggleSidebar = document.querySelector("#toggleSidenav");
let sidebar = document.querySelector(".side-header");
let sideNavLinks = document.querySelectorAll(".side-nav a");
let main = document.querySelector("main");
let currentActive = document.querySelector(".side-nav .active");
let currentPageName = document.querySelector(".current-page-name");
let currentScript = localStorage.getItem("currentScript");
let content = document.querySelector("#content");




window.addEventListener("load", () => {
	document.body.style.overflow = "hidden";
	setTimeout(() => {
		preloader.animate([{
			opacity: "1"
		}, {
			opacity: "0",
			display: "none"
		}], {
			duration: 1000,
			fill: "forwards",
		});
		document.body.style.overflow = "auto";
	}, 1000);
	preloader.addEventListener("animationend", () => {
		preloader.style.display = "none";


	});

	let index = localStorage.getItem("currentPageIndex")
	setTimeout(() => {
		addScript(currentScript);
	}, 1000);
	sideNavLinks[index].click();
	
});

toggleSidebar.addEventListener("click", () => {
	sidebar.classList.toggle("collapse");
	if (sidebar.classList.contains("collapse")) {
		main.style.margin = "0 0 0 3vw";
		sideNavLinks.forEach((link) => {
			link.style.justifyContent = "center";
		});
	} else {
		main.style.margin = "0 0 0 13vw";
		sideNavLinks.forEach((link) => {
			link.style.justifyContent = "start";
		});
	}
});

sideNavLinks.forEach((link,index) => {
	link.addEventListener("click", () => {
		let htmlFile = link.getAttribute("data-link");
		let scriptSrc = link.getAttribute("data-script")		
			setTimeout(() => {
				if (scriptSrc !== currentScript) {
					deleteScript(currentScript)
					addScript(scriptSrc);
					currentScript = scriptSrc;
					localStorage.setItem("currentScript",scriptSrc)
				}
	}, 1000);
		let title = link.textContent;
		currentActive.classList.remove("active");
		link.parentElement.classList.add("active");
		currentActive = link.parentElement;
		localStorage.setItem("currentPageIndex",index)
		currentPageName.textContent = title;
		showContent(htmlFile, title);
	});
});


function showContent(htmlFile, title) {
	fetch(htmlFile)
		.then(response => {
			if (!response.ok) {
				throw new Error("HTTP error " + response.status);
			}
			return response.text();
		})
		.then(data => {
			content.innerHTML = data;
			document.title = "DEKKOTO - " + title;
		})
		.catch(function () {
			console.log('Fetch Error :-S', err);
		});
}


// Attach an event listener to the document body



document.addEventListener('click', function (e) {
	let fileInput;
	let selectName;
	let imgPreview;
	// console.log(e.target)

	if (e.target.matches('#selectFile')) {
		console.log('selectFile')
		fileInput = document.getElementById('uploadVideo');
		fileInput.click();
		selectName = document.getElementsByClassName("file-select-name")[0];
	} else if (e.target.matches('#selectThumbnail')) {
		fileInput = document.getElementById('uploadThumbnail');
		fileInput.click();
		imgPreview = document.querySelector('.previewImgThumbnail');
	} else if (e.target.matches('#selectBanner')) {
		fileInput = document.getElementById('uploadBanner');
		fileInput.click();
		imgPreview = document.querySelector('.previewImgBanner');
	}

	if (fileInput) {
		fileInput.addEventListener("change", function () {
			if (selectName) {
				let filename = fileInput.files[0].name;
				selectName.innerText = filename;
			} else if (imgPreview) {
				let file = fileInput.files[0];
				imgPreview.src = URL.createObjectURL(file);
			}
		});
	}
});

content.addEventListener('drop', function (e) {
	e.preventDefault();
	e.stopPropagation();
	let fileInput = document.getElementById('uploadVideo');
	if (e.target.closest('.upload-area')) {
		let dt = new DataTransfer();

		for (let file of e.dataTransfer.files) {
			dt.items.add(file);
		}
		fileInput.files = dt.files;
		let filename = fileInput.files[0].name;
		let selectName = document.getElementsByClassName("file-select-name")[0];
		selectName.innerText = filename;
	}
});

content.addEventListener('dragover', function (e) {
	e.preventDefault();
	e.stopPropagation();
	if (e.target.matches('.upload-area')) {
		document.querySelector('.upload-area').style.backgroundColor = 'rgba(0, 0, 0, 0.1)';
	}
});

content.addEventListener('dragleave', function (e) {
	e.preventDefault();
	e.stopPropagation();
	if (e.target.matches('.upload-area')) {
		document.querySelector('.upload-area').style.backgroundColor = '';
	}
});

function addScript(src) {
		let script = document.createElement("script");
		script.src = src;
		script.defer = true;
		document.body.appendChild(script);
}

function deleteScript(src) {
	let script = document.querySelectorAll("script")
	script.forEach((js) => {
		let scriptUrl = new URL(js.src);
		let pathname = scriptUrl.pathname;
		let scriptSrc = ".." + pathname
		if (scriptSrc === src) {
			document.body.removeChild(js);
		}
	});
}

