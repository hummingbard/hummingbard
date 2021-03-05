<script>
import Avatar from '../../../avatar/avatar.svelte'
import { settings,state } from '../store.js'


function updateAvatar(e) {
  console.log("lol", e.detail)
  $settings.info.avatar = e.detail
}

function removeAvatar(e) {
  $settings.info.avatar = null
}


$: about = window.timeline?.profile ? `About you` : `About this community`

</script>


<div class="flex flex-column flex-one ">
  <div class="flex flex-column ">
    <div class="">
      <span class="small"><strong>Title</strong></span>
    </div>
    <div class="mt2">
      <input
        placeholder="Title"
        bind:value={$settings.info.title}
      />
    </div>
  </div>


  <div class="mt3 flex flex-column">
    <div class="">
      <span class="small"><strong>About</strong></span>
    </div>
    <div class="mt2">
      <textarea
        style="height:100px;"
        placeholder={about}
        bind:value={$settings.info.about}
      ></textarea>
    </div>
  </div>
</div>

<div class="ml4 flex flex-column">

  <div class="">
    <Avatar 
      on:uploaded={updateAvatar}
      on:removed={removeAvatar}
      avatar={$settings.info.avatar}/>
  </div>

  <div class="">
  </div>

</div>
