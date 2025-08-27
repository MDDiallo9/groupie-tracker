let rangeMin = 1;
const range = document.querySelector(".range-selected");
const rangeInput = document.querySelectorAll(".range-input input");
const rangePrice = document.querySelectorAll(".range-price input");

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
