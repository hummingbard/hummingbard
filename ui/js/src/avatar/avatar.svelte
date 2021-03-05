<script>
import {onMount, createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher()

export let avatar;
export let header;

$: exists = avatar?.length > 0

let imageInput;

function add() {
    imageInput.click()
}

let uploaded = false;
let uploading = false;
let file;
let url;

let mxc;

$: imgSrc = uploading ? url : exists ? `${homeserverURL}/_matrix/media/r0/thumbnail/${avatar?.substring(6)}?width=100&height=100&method=scale` : ``

let build = (e) => {


    const allowed = ["image/jpeg", "image/jpg", "image/webp", "image/png"]

    file = e.target.files[0]

    if(!file) {
        return
    }

    if (file && !allowed.includes(file.type)) {
        alert("That is not a valid image.")
        return
    }
    if (file?.size > 23000000) {
        alert("That image is too large.")
        return
    }


    var reader = new FileReader();
    reader.readAsDataURL(file);

    reader.onload = e => {
      const content = e.target.result;

        url = URL.createObjectURL(file);

        uploading = true

          uploadImage().then((res) => {
              if(res?.content_uri) {
                  uploaded = true
                  mxc = res.content_uri
              }
          }).then(() => {
              uploading = false
              uploaded = true
              dispatch('uploaded', mxc)
          })

    }


}

async function uploadImage() {
    let endpoint = `${identity.well_known}/_matrix/media/r0/upload`

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: file,
    headers:{
        'Authorization': `Bearer ${identity.matrix_access_token}`,
        'Content-Type': file.type
    }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


function remove() {
    avatar = null
    file = null
    url = null
    uploaded = false
  dispatch('removed', true)
}

</script>


<input 
    type="file" 
    accept="image/jpeg, image/png, image/webp" 
    name="images"
    bind:this={imageInput} 
    on:change={build} 
    hidden 
>

<div class="flex flex-column">

    <div class="avatar gr-default pointer" 
         class:header={header}
         class:o-50={uploading}
         style="background-image:url({imgSrc})"
        on:click={add}>

        {#if !exists}
            <div class="gr-center">
                <svg class="sv" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M19.25 4.5H4.75a.25.25 0 00-.25.25v14.5c0 .138.112.25.25.25h.19l9.823-9.823a1.75 1.75 0 012.475 0l2.262 2.262V4.75a.25.25 0 00-.25-.25zm.25 9.56l-3.323-3.323a.25.25 0 00-.354 0L7.061 19.5H19.25a.25.25 0 00.25-.25v-5.19zM4.75 3A1.75 1.75 0 003 4.75v14.5c0 .966.784 1.75 1.75 1.75h14.5A1.75 1.75 0 0021 19.25V4.75A1.75 1.75 0 0019.25 3H4.75zM8.5 9.5a1 1 0 100-2 1 1 0 000 2zm0 1.5a2.5 2.5 0 100-5 2.5 2.5 0 000 5z"></path></svg>
            </div>
        {/if}

    </div>

    <div class="mt1" class:tc={!header}>
        {#if uploaded || exists}
            <span class="small pointer" on:click={remove}><u>Remove</u></span>
        {/if}
    </div>
</div>


<style>
.avatar {
    background-color: var(--primary-darkest);
    border-radius: 50%;
    width: 100px;
    height: 100px;
    background-repeat: no-repeat;
    background-size: cover;
    background-position: center;
}

.header {
    border-radius: 7px;
    width: 100%;
}

.avatar:hover {
    opacity: 0.9;
}

.sv {
    fill: white;
}
</style>
