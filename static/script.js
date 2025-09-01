let rangeMin = 1;
const range = document.querySelector(".range-selected");
const rangeInput = document.querySelectorAll(".range-input input");
const rangePrice = document.querySelectorAll(".range-price input");

const listElements = document.querySelectorAll('li');


const searchForm = document.querySelector('#searchBar');
searchForm.addEventListener("submit", searchSubmit)

async function searchSubmit(e) {
    e.preventDefault()
    const response = await fetch('http://localhost:8000/search', {
        method: 'POST',
        body: searchForm[0].value
    });
    const data = await response.json();
    const main = document.querySelector('main');
    while (main.children.length > 1) {
        main.removeChild(main.lastChild);
    }
    const sep = document.createElement("div")
    sep.classList.add("separator")
    sep.textContent = `Résultats pour ${searchForm[0].value}`
    main.appendChild(sep)
    const result = document.createElement("div")
    result.classList.add("result")
    data.forEach(artist => {
        const slide = document.createElement('a');
        slide.className = 'card';
        slide.href = `artist?id=${artist.id}`;
        slide.innerHTML = `
      <div class="artistImage">
        <img src="${artist.image}" alt="">
      </div>
      <p class="artistName">${artist.name}</p>
    `;
        result.appendChild(slide);
    });
    main.appendChild(result)
    searchForm[0].value = ""
}



rangeInput.forEach((input) => {
    input.addEventListener("input", (e) => {
        let minRange = parseInt(rangeInput[0].value);
        let maxRange = parseInt(rangeInput[1].value);
        if (maxRange - minRange < rangeMin) {
            if (e.target.className === "min") {
                rangeInput[0].value = maxRange - rangeMin;
            } else {
                rangeInput[1].value = minRange + rangeMin;
            }
        } else {
            rangePrice[0].value = minRange;
            rangePrice[1].value = maxRange;
            range.style.left = (minRange / rangeInput[0].max) * 100 + "%";
            range.style.right = 100 - (maxRange / rangeInput[1].max) * 100 + "%";
        }
    });
});

rangePrice.forEach((input) => {
    input.addEventListener("input", (e) => {
        let minPrice = rangePrice[0].value;
        let maxPrice = rangePrice[1].value;
        if (maxPrice - minPrice >= rangeMin && maxPrice <= rangeInput[1].max) {
            if (e.target.className === "min") {
                rangeInput[0].value = minPrice;
                range.style.left = (minPrice / rangeInput[0].max) * 100 + "%";
            } else {
                rangeInput[1].value = maxPrice;
                range.style.right = 100 - (maxPrice / rangeInput[1].max) * 100 + "%";
            }
        }
    });
});

// First Album Date
const minfAdRange = document.getElementById("minfAdRange");
const maxfAdRange = document.getElementById("maxfAdRange");
const minfAd = document.getElementById("minfAd");
const maxfAd = document.getElementById("maxfAd");
const minGapFAD = 1;

minfAdRange.addEventListener("input", function () {
    let minVal = parseInt(minfAdRange.value);
    let maxVal = parseInt(maxfAdRange.value);
    if (minVal > maxVal - minGapFAD) {
        minVal = maxVal - minGapFAD;
        minfAdRange.value = minVal;
    }
    minfAd.value = minVal;
});
maxfAdRange.addEventListener("input", function () {
    let minVal = parseInt(minfAdRange.value);
    let maxVal = parseInt(maxfAdRange.value);
    if (maxVal < minVal + minGapFAD) {
        maxVal = minVal + minGapFAD;
        maxfAdRange.value = maxVal;
    }
    maxfAd.value = maxVal;
});
minfAd.addEventListener("input", function () {
    let minVal = parseInt(minfAd.value);
    let maxVal = parseInt(maxfAd.value);
    if (minVal > maxVal - minGapFAD) {
        minVal = maxVal - minGapFAD;
        minfAd.value = minVal;
    }
    minfAdRange.value = minVal;
});
maxfAd.addEventListener("input", function () {
    let minVal = parseInt(minfAd.value);
    let maxVal = parseInt(maxfAd.value);
    if (maxVal < minVal + minGapFAD) {
        maxVal = minVal + minGapFAD;
        maxfAd.value = maxVal;
    }
    maxfAdRange.value = maxVal;
});

// Creation Date
const minCDRange = document.getElementById("minCDRange");
const maxCDRange = document.getElementById("maxCDRange");
const minCD = document.getElementById("minCD");
const maxCD = document.getElementById("maxCD");
const minGapCD = 1;

minCDRange.addEventListener("input", function () {
    let minVal = parseInt(minCDRange.value);
    let maxVal = parseInt(maxCDRange.value);
    if (minVal > maxVal - minGapCD) {
        minVal = maxVal - minGapCD;
        minCDRange.value = minVal;
    }
    minCD.value = minVal;
});
maxCDRange.addEventListener("input", function () {
    let minVal = parseInt(minCDRange.value);
    let maxVal = parseInt(maxCDRange.value);
    if (maxVal < minVal + minGapCD) {
        maxVal = minVal + minGapCD;
        maxCDRange.value = maxVal;
    }
    maxCD.value = maxVal;
});
minCD.addEventListener("input", function () {
    let minVal = parseInt(minCD.value);
    let maxVal = parseInt(maxCD.value);
    if (minVal > maxVal - minGapCD) {
        minVal = maxVal - minGapCD;
        minCD.value = minVal;
    }
    minCDRange.value = minVal;
});
maxCD.addEventListener("input", function () {
    let minVal = parseInt(minCD.value);
    let maxVal = parseInt(maxCD.value);
    if (maxVal < minVal + minGapCD) {
        maxVal = minVal + minGapCD;
        maxCD.value = maxVal;
    }
    maxCDRange.value = maxVal;
});

const swiper = new Swiper('.playlist-swiper', {
    slidesPerView: 'auto',
    spaceBetween: 24,
    grabCursor: true,
    navigation: {
        nextEl: '.swiper-button-next',
        prevEl: '.swiper-button-prev',
    },
    effect: 'slide',
});

listElements.forEach(list => {
    list.addEventListener("click", () => {
        console.log("clicked")
        refreshPlaylistSwiper(list, swiper);
    });
});

// Popular swiper

const swiperP = new Swiper('.popular-swiper', {
    effect: 'coverflow',
    grabCursor: true,
    centeredSlides: true,
    slidesPerView: 'auto',
    loop: true,
    coverflowEffect: {
        rotate: 0,
        stretch: 0,
        depth: 200,
        modifier: 1,
        slideShadows: true,
    },
    navigation: {
        nextEl: '.swiper-button-next-popular',
        prevEl: '.swiper-button-prev-popular',
    },
    pagination: {
        el: '.swiper-pagination-popular',
        clickable: true,
    },
});

// Color thief

const colorThief = new ColorThief();
const img = document.querySelector('#highlightImage');
const image = document.querySelector('#highlightImage');
img.crossOrigin = "anonymous";

// Make sure image is finished loading
if (img.complete) {
    console.log(colorThief.getColor(img));
} else {
    image.addEventListener('load', function () {
        console.log(colorThief.getColor(img));
    });
}

// Spotify API


const highlightName = document.querySelector('#highlightName');


// Refresh Playlist

async function refreshPlaylistSwiper(element, swiper) {
    element.parentElement.querySelectorAll('li').forEach(li => {
        li.classList.remove('active');
    });

    element.classList.add('active');
    const res = await fetch('http://localhost:8000/api?' + element.dataset.api);
    const artists = await res.json();
    console.log(artists)

    let wrapper;
    if (element.dataset.api.includes('dec')) {
        wrapper = document.querySelector('#decWrapper.swiper-wrapper');
    } else {
        wrapper = document.querySelector('#locWrapper.swiper-wrapper');
    }
    wrapper.textContent = ''
    artists.forEach(artist => {
        const slide = document.createElement('a');
        slide.className = 'swiper-slide';
        slide.href = `artist?id=${artist.id}`;
        slide.innerHTML = `
      <div class="artistImage">
        <img src="${artist.image}" alt="">
      </div>
      <p class="artistName">${artist.name}</p>
    `;
        wrapper.appendChild(slide);
    });
    swiper.forEach(s => {
        s.update()
    })

}




