console.log(`%cHummingbard`, 'color: orange; font-size: 23px;font-weight:bold')

// SET UP back-to-top button
let observer;

function callback(entries, observer) {
  entries.forEach(entry => {
    if(!entry.isIntersecting) {
      sttp.classList.remove('dis-no')
    } else {
      sttp.classList.add('dis-no')
    }
  })
}

let options = {
    root: document.body,
    rootMargin: `0px`,
    threshold: 1,
}

let sttp = document.querySelector('.sttp')
let header = document.querySelector('.header')

if(sttp && header) {
  sttp.addEventListener('click', () => {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  })
  observer = new IntersectionObserver(callback);
  observer.observe(header)
}

let dm = document.querySelector('.dark-mode')
if(dm) {
  import('./dark-mode/dark-mode.svelte').then(res => {
      new res.default({
          target: dm,
          hydrate: true
      });
  })
}
