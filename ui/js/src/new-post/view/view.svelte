<script>
import { onMount, tick, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher();
import {fade,fly} from 'svelte/transition'
import { text_area_resize } from '../../utils/auto-resize-textarea.js'

import { addedPosts, addPost } from '../../timeline/store.js'
import { post, articleSettings } from '../store.js'
import {makeid} from '../../utils/utils.js'
import ItemsView from './items-view.svelte'
import Tools from './tools/tools.svelte'
import {editorState} from '../../editor/store.js'

export let reply = false;
export let eventID;

onMount(() => {
  if(reply){
    store.active = true
  }
})

// let's store props here for now until i figure out a better way
let store = {
  active: false,
  content: {
    plain_text: null,
    html: null,
  },
  links: [],
  attachments: [],
  images: [],
  locked: false,
  nsfw: false,
  anonymous: false,
  article: {
    enabled: false,
    title: null,
    subtitle: null,
    description: null,
    canonical_link: null,
    featured_image: null,
    settings: {
      active: false,
    }
  },
}

function kill() {
  store.nsfw = false
  store.anonymous = false
  store.article.enabled = false
  store.article.title = null
  store.article.subtitle = null
  store.article.description = null
  store.article.canonical_link = null
  store.article.featured_image = null
  post.kill()
}

function toggleArticleOn() {
  store.article.enabled = true
}

function toggleArticleOff() {
  store.article.enabled = false
}

function killArticleSettings() {
  store.article.settings.active = false
}

function toggleArticleSettings() {
  store.article.settings.active = !store.article.settings.active
}

function updateFeaturedImage(e) {
  store.article.featured_image = e.detail
}

function removeFeaturedImage(e) {
  store.article.featured_image = null
}

function updateDetails(e) {
  store.article.subtitle = e.detail.subtitle
  store.article.description = e.detail.description
  store.article.canonical_link = e.detail.canonical_link
}


function toggleAnonymous() {
  store.anonymous = !store.anonymous
}
function toggleNSFW() {
  store.nsfw = !store.nsfw
}
function lock() {
  store.locked = true
}
function unlock() {
  store.locked = false
}

function addLink(e) {
  let ind = store.links.findIndex(x => x.href === e.detail.href)
  if(ind == -1) {
    store.links.push(e.detail)
    store.links = store.links
  }
}
function updateLinkMetadata(e) {
  let link = e.detail
  let ind = store.links.findIndex(x => x.id === link.id)
  if(ind != -1) {
    store.links[ind].metadata = link.metadata
  }
}

function deleteLink(e) {
  let id = e.detail
  let ind = store.links.findIndex(x => x.id === id)
  if(ind != -1) {
    store.links.splice(ind, 1)
    store.links = store.links
  }
}

function addImage(e) {
  let ind = store.images.findIndex(x => x.id === e.detail.id)
  if(ind == -1) {
    store.images.push(e.detail)
    store.images = store.images
  }
}

function deleteImage(e) {
  let ind = store.images.findIndex(x => x.id === e.detail)
  if(ind != -1) {
    store.images.splice(ind, 1)
    store.images = store.images
  }
}

function updateImageURL(e) {
  let ind = store.images.findIndex(x => x.id === e.detail.id)
  if(ind != -1) {
    store.images[ind].mxc = e.detail.content_uri
  }
}

function updateImageMetadata(e) {
  let metadata = e.detail
  let ind = store.images.findIndex(x => x.id === metadata.id)
  if(ind != -1) {
    store.images[ind].caption = metadata.metadata.caption
    store.images[ind].description = metadata.metadata.description
  }
}

function addAttachment(e) {
  let ind = store.attachments.findIndex(x => x.id === e.detail.id)
  if(ind == -1) {
    store.attachments.push(e.detail)
    store.attachments = store.attachments
  }
}

function deleteAttachment(e) {
  let ind = store.attachments.findIndex(x => x.id === e.detail)
  if(ind != -1) {
    store.attachments.splice(ind, 1)
    store.attachments = store.attachments
  }
}

function updateAttachmentURL(e) {
  let ind = store.attachments.findIndex(x => x.id === e.detail.id)
  if(ind != -1) {
    store.attachments[ind].mxc = e.detail.content_uri
  }
  console.log(store)
}


let editor;

$: active = reply? store.active : $post.active

let Editor;
let editorLoaded;
$: if(active) {
  import('../../editor/editor.svelte').then(res => {
    Editor = res.default
    editorLoaded = true
  })
}


let view;
async function scrollIntoView() {
  view.scrollIntoView({ behavior: 'smooth', block: 'start'});
}

/*
$: if(active && editorLoaded && window.timeline?.permalink && view && !reply) {
  scrollIntoView()
}
*/


$: empty = (!store.content?.plain_text || store.content?.length == 0) &&
  store.images.length == 0 &&
  store.links.length == 0 &&
  store.attachments.length == 0 

$: locked = store.locked

function sync(e) {
  let content = e.detail
  store.content = {
    plain_text: content.plain_text,
    html: content.html,
    length: content.length,
  }
}


function listener() {
    let el = document.querySelector('.editor-content')

    el.addEventListener('paste', (e) => {
        e.preventDefault()
      if(store.article.enabled) {
        return
      }
            let expression = /(http|ftp|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])?/g
            let regex = new RegExp(expression);
     
            let paste = (event.clipboardData || window.clipboardData).getData('text');
            let matches = paste.match(regex);
            if(matches && matches.length > 0) {
                for(let i=0;i<matches.length; i++) {
                    if(i == 9) {
                        break
                    }
                    pastedLink(matches[i])
                }
            }
    }); 

}

function doesLinkExist(href) {
    let ind = store.links.findIndex(x => x.href === href)
    if(ind != -1) {
        return true
    }
    return false
}

function pastedLink(href) {
    const exists = doesLinkExist(href)
    if(!exists) {
        let item = {
            id: makeid(12),
            href: href,
            metadata: {
                title: null,
                description: null,
                author: null,
                image: null,
                domain: null,
            },
            data: null,
        }
      console.log(item)
        store.links.push(item)
        store.links = store.links
    }
}


$: text = window?.timeline?.permalink ? `Reply` : `Post`

async function focusTitle() {
  await tick;
}


function create() {
  if(empty) {
    alert("Post can't be empty.")
    return
  }
  if(isArticle && (!titleInput.value || titleInput.value.length == 0)) {
    alert("Your article must have a title.")
    focusTitle()
    return
  }
  lock()
  createPost().then((res) => {
    console.log(res)
    if(res?.post) {
      res.post.just_added = true
      if(!store.anonymous) {
        res.post.author.display_name = identity.display_name
        res.post.author.avatar_url = identity.avatar_url
      }
      res.post.replies = []
      if(reply) {
        dispatch('added', {
          id: eventID,
          post: res.post,
        })
      } else {
        addPost(res.post)
        addedPosts.add(res.post)
      }
      discard()
    } else {
      unlock()
    }
  }).then(() => {
      unlock()
    if(window.timeline?.permalink) {
      let el = document.querySelector('.no-replies')
      if(el) {
        el.remove()
      }
    }
  })
}

async function createPost() {
    let endpoint = `/post/create`

    let data = {
        room_id: window.timeline.room_id,
        room_alias: window.timeline.alias,
        room_path: window.timeline.room_path,
        post: {
          content: {
              text: store.content.plain_text,
          },
          article: {
            enabled: false,
          }
        },
      nsfw: store.nsfw,
      anonymous: store.anonymous,
    }

  if(window.timeline?.userFeed && window.timeline?.feed) {
    data.room_id = identity?.room_id
    data.room_alias = `#${identity.user_id}`
  }

  if(window?.timeline?.permalink && window?.timeline?.event_id) {
    //data.room_id = window.timeline?.permalinkedPost?.content?.thread_room_id
    data.room_alias = window.timeline?.permalinkedPost?.content?.thread_room_alias
    //data.room_id = window.timeline.thread_in_room_id
    data.reply = true
    data.event_id = window.timeline.event_id
    if(reply) {
      data.event_id = eventID
    }
  }

  if(store.links.length > 0) {
    data.post.links = store.links
  }

  if(store.images.length > 0) {
    let images = [];
    store.images.forEach(image => {
      images.push({
        caption: image.caption,
        description: image.description,
        filename: image.file.name,
        size: image.file.size,
        mimetype: image.file.type,
        mxc: image.mxc,
        width: image.width,
        height: image.height,
      })
    })
    data.post.images = images
  }

  if(store.attachments.length > 0) {
    let attachments = [];
    store.attachments.forEach(attachment => {
      attachments.push({
        filename: attachment.file.name,
        size: attachment.file.size,
        mimetype: attachment.file.type,
        mxc: attachment.mxc,
      })
    })
    data.post.attachments = attachments
  }

  if(store.article.enabled) {
    data.post.article = {
      enabled: true,
      title: titleInput.value,
      subtitle: store.article.subtitle,
      description: store.article.description,
      canonical_link: store.article.canonical_link,
    }
    if(store.article.featured_image) {
      data.post.article.featured_image = {
        caption: store.article.featured_image.caption,
        mxc: store.article.featured_image.content_uri,
        height: store.article.featured_image.height,
        width: store.article.featured_image.width,
      }
    }
  }

  console.log(data)

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: JSON.stringify(data),
    headers:{
      'Authorization': identity.access_token,
        'Content-Type': 'application/json'
    }
    })
    if (resp.ok) { // if HTTP-status is 200-299
    } else {
      alert("HTTP-Error: " + resp.status);
    }
    const ret = await resp.json()
    return Promise.resolve(ret)
}

function discard() {
  killArticleSettings()
  kill()
  if(reply) {
    dispatch('killed')
  }
}

$: isArticle = store.article?.enabled

function toggleArticle() {
  post.toggleArticle()
}

let mdLoaded;
let Markdown;
$: if(isArticle) {
  killBodyScroll()
  import('./tools/article/markdown.svelte').then(res => {
    Markdown = res.default
    mdLoaded = true
  })
} else {
  unkillBodyScroll()
}

function killBodyScroll() {
  let nav = document.querySelector('.nav-de') 
  if(nav) {
    nav.classList.add('hide')
  }
}

function unkillBodyScroll() {
  let nav = document.querySelector('.nav-de') 
  if(nav) {
    nav.classList.remove('hide')
  }
}


let titleInput;
function updateTitle(e) {
  if(e.key == 'Enter') {
    e.preventDefault()
    return
  }
}

$: settings = store?.article?.settings?.active

let Settings;
let settingsLoaded = false;
function loadSettings() {
  import('../view/settings/settings.svelte').then(res => {
    Settings = res.default
    settingsLoaded = true
  })
}

$: if(settings) {
  loadSettings()
}

$: featuredImageExists = store?.article?.featured_image?.content_uri?.length > 0
$: featuredImage = `${homeserverURL}/_matrix/media/r0/thumbnail/${store?.article?.featured_image?.content_uri?.substring(6)}?width=100&height=100&method=scale`

$: allowArticle = window.timeline?.permalink ? window.timeline.permalinkedPost.content['com.hummingbard.article'] ? true : false : true

</script>


{#if active && editorLoaded}

<div class="co grr w-100" 
bind:this={view}
class:pia={!isArticle}
class:mb3={!reply}
class:reply={reply}
class:article={isArticle}
class:brd-a={window.timeline?.is_article}
transition:fade="{{duration: 33}}">
  <div class="flex flex-column grr" 
    class:pa3={!isArticle}
    class:ovfl-y={isArticle}
    class:scrl={isArticle}>

    <div class="contain relative">


        {#if settings && settingsLoaded}
          <div class="article-settings grr"
          transition:fly="{{ x: -350, duration: 80 }}">
            <Settings 
              on:updateFeaturedImage={updateFeaturedImage}
              on:removeFeaturedImage={removeFeaturedImage}
              on:updateDetails={updateDetails}
              store={store}/>
            <div class="">
            </div>
          </div>
        {/if}


      <div class="flex flex-column relative">
        {#if isArticle}

          {#if mdLoaded}
            <Markdown/>
          {/if}

          <div class="article-title">
            <textarea 
              class="title"
              bind:this={titleInput}
              on:keypress={updateTitle}
              on:paste|preventDefault=""
              placeholder="Title"
              autofocus
              use:text_area_resize/>
          </div>
        {/if}

        {#if featuredImageExists}
          <div class="e-pad-a pt2 pb3 flfex flex-column">
            <div class="relative">
              <img src={featuredImage} />
              <div class="discard pointer o-70 hov-op"
                on:click={removeFeaturedImage}>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M1 12C1 5.925 5.925 1 12 1s11 4.925 11 11-4.925 11-11 11S1 18.075 1 12zm8.036-4.024a.75.75 0 00-1.06 1.06L10.939 12l-2.963 2.963a.75.75 0 101.06 1.06L12 13.06l2.963 2.964a.75.75 0 001.061-1.06L13.061 12l2.963-2.964a.75.75 0 10-1.06-1.06L12 10.939 9.036 7.976z"></path></svg>
              </div>
            </div>
            <div class="mb2">
              <input 
                class="caption"
                bind:value={store.article.featured_image.caption} 
                placeholder="Caption"/>
            </div>


          </div>
        {/if}

        <div class="h-100">
          <Editor 
            article={isArticle}
            bind:this={editor} 
            on:ready={listener}
            on:sync={sync}/>
        </div>
        {#if !isArticle}
        <div class="">
          <ItemsView
            store={store} 
            on:deleteLink={deleteLink}
            on:updateLinkMetadata={updateLinkMetadata}
            on:deleteImage={deleteImage}
            on:updateImageURL={updateImageURL}
            on:updateImageMetadata={updateImageMetadata}
            on:deleteAttachment={deleteAttachment}
            on:updateAttachmentURL={updateAttachmentURL}
            />
        </div>
        {/if}


      </div>

    </div>

  </div>
  <div class={isArticle ? 'pa3 flex tia tools' : 'ph3 pb3 flex'}>
      <div class="flex-one gr-center">
        <Tools 
          store={store} 
          allowArticle={allowArticle} 
          on:toggleArticleOn={toggleArticleOn}
          on:toggleArticleOff={toggleArticleOff}
          on:killArticleSettings={killArticleSettings}
          on:addLink={addLink}
          on:addImage={addImage}
          on:addAttachment={addAttachment}
          on:toggleNSFW={toggleNSFW}
          on:toggleAnonymous={toggleAnonymous}
          on:toggleArticleSettings={toggleArticleSettings}/>
      </div>
      {#if store.locked}
        <div class="mh3 gr-center">
          <div class="lds-ring"><div></div><div></div><div></div><div></div></div>
        </div>
      {/if}
      <div class="">
        <button class="" on:click={create} disabled={store.locked}>{text}</button>
        <button class="light" on:click={discard}>Cancel</button>
      </div>
    </div>
</div>
{/if}


<style>

.co {
  min-height: 180px;
}

.grr {
  display: grid;
  grid-template-rows: [editor] 1fr [tools] auto;
}

.contain {
  display: grid;
  grid-template-rows: auto;
  grid-template-columns: auto;
}

.contain-s {
  grid-template-columns: [settings] 350px [editor] 1fr;
}

.article {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 999;
  background: var(--m-bg);
}

.brd-a {
  border-left: 1px solid var(--primary-grayish);
  border-right: 1px solid var(--primary-grayish);
  border-top: 1px solid var(--primary-grayish);
}

.article-title {
}

.article-settings {
  position: fixed;
  left: 0;
  top:0;
  bottom: 0;
  width: 350px;
  border-right: 1px solid var(--primary-grayish);
  z-index: 1000;
  background: var(--m-bg);
    box-shadow: 0 30px 60px rgba(0,0,0,.05);
}

.tools {
  z-index: 1001;
  background: var(--m-bg);
}


.title {
  font-size: 1.5rem;
  font-weight: bold;
  border: 0;
  padding: 0;
  margin: 0;
    --minus: calc(100% - 720px);
    --bytwo: calc(var(--minus) / 2);
    padding-left: var(--bytwo);
    padding-right: var(--bytwo);
    padding-top: 2rem;
    padding-bottom: 0.75rem;
}

@media screen and (max-width: 768px) {
  .title {
    padding-left: 20px;
    padding-right: 20px;
  }
}

@media screen and (max-width: 568px) {
  .article-settings {
    width: 100%;
  }
}


.pia {
    background-color: var(--pi-bg);
    transition: 0.1s;
    word-break: break-word;
    border-bottom: 1px solid var(--primary-grayish);
}

.tia {
  border-top: 1px solid var(--primary-grayish);
}


.caption {
  border: 0;
  text-align: center;
  font-size: small;
}

.discard {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
}
.reply {
  border-bottom: 1px solid var(--primary-grayish);
}
  textarea, input {
      width: 100%;
  }
</style>
