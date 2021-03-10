<script>
import { createEventDispatcher } from 'svelte';
const dispatch = createEventDispatcher();

import {makeid} from '../../../../utils/utils.js'

export let store;

let imageInput;

let files = [];
let count = 0;


export let imageAllowed = true;


function add() {
    if(!imageAllowed) {
        return
    }
    if(store.images.length >= 14) {
        return
    }
    imageInput.click()
}


let build = (e) => {
    if(count + e.target.files.length > 14) {
        alert("That's too many images for one post.")
        return
    }
    for(let i =0 ; i < e.target.files.length ; i++) {

        if(e.target.files.length > 14) {
            alert("That's too many images for one post.")
            break
        }

        const allowed = ["image/gif","image/jpeg", "image/jpg", "image/webp", "image/png"]

        const file = e.target.files[i]

        if (file && !allowed.includes(file.type)) {
            alert("That is not a valid image.")
          continue
        }
        if (file.size > 23000000) {
            alert("That image is too large.")
            continue
        }
        files = [...files, e.target.files[i]]



    }

    count += files.length


    for(let i =0 ; i < files.length; i++) {

        var reader = new FileReader();
        const file = files[i]
        reader.readAsDataURL(file);

        reader.onload = e => {
          const content = e.target.result;

          let item = {
              id: `im-${makeid(16)}`,
              url: URL.createObjectURL(file),
              file: file,
              caption: '',
              description: '',
              uploaded: false,
          }

          var image = new Image();
          image.src = item.url
          image.onload = () => {
            item.height = image.height
            item.width = image.width
              dispatch('addImage', item)
              imageInput.value = ''
          }


        }
    }

    files = []



}



</script>

<input 
    type="file" 
    accept="image/gif, image/jpeg, image/png, image/webp" 
    name="images"
    bind:this={imageInput} 
    on:change={build} 
    hidden 
    multiple
>
<div class="pointer o-70 hov-op" class:fill-blue={store?.images?.length > 0} 
  aria-label="Add Image"
  data-microtip-position="bottom"
  data-microtip-size="fit"
  role="tooltip"
    on:click={add}>
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M1.75 2.5a.25.25 0 00-.25.25v10.5c0 .138.112.25.25.25h.94a.76.76 0 01.03-.03l6.077-6.078a1.75 1.75 0 012.412-.06L14.5 10.31V2.75a.25.25 0 00-.25-.25H1.75zm12.5 11H4.81l5.048-5.047a.25.25 0 01.344-.009l4.298 3.889v.917a.25.25 0 01-.25.25zm1.75-.25V2.75A1.75 1.75 0 0014.25 1H1.75A1.75 1.75 0 000 2.75v10.5C0 14.216.784 15 1.75 15h12.5A1.75 1.75 0 0016 13.25zM5.5 6a.5.5 0 11-1 0 .5.5 0 011 0zM7 6a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
</div>

