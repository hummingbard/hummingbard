<script>
import { tick} from 'svelte'
import {fade} from 'svelte/transition'
import Image from './image/image.svelte'
import Link from './link/link.svelte'
import Attachment from './attachment/attachment.svelte'
import Menu from './menu/menu.svelte'
import Content from './content/content.svelte'
import Reactions from './reactions/reactions.svelte'
import {addPosts} from '../timeline/store.js'
import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher();

export let post;
export let embed = false;;
export let nested = false;;
export let reply = false;;
export let nestedReply = false;;
export let depth = 0;



export let expanded = false;;

let expand = expanded;

function toggleExpand(e) {
    if(expanded) {
        dispatch('collapse')
        return
    }
    expand = !expand
    if(!expand) {
        collapse()
    }
}

function expandReplies(e) {
    expand = true
}
function collapseReplies(e) {
    expand = false
}


function collapse(e) {
    expand = false
    dispatch('collapse')
}

$: senderHREF = `/${post.author.formatted_id}`

$: user = post.author?.display_name?.length > 0 ? post.author?.display_name : post.author?.formatted_id

$: showPath = (post?.content?.room_path != window?.timeline?.room_path) &&
    !post?.content?.room_path?.includes('@')

$: rpth = (window.timeline?.room_path == post.content.room_path) ||
    window.timeline?.room_path == 'user-index'

$: defaultLink = rpth ? `/${post.content.room_path}/${post.event_id}` :
    nested ? `${post.content.room_path}/${post.event_id}`:`/${post.content.room_path}/${post.content.event_id}`

$: direct = window.timeline?.permalink ? `/${window.timeline.root_event}/${post.event_id}` : defaultLink


$: permalink = (nested && post?.content?.share_reply) ?  `/${post.content?.reply_permalink}` : direct

$: mediaExists = post?.content?.images?.length > 0 ||
    post?.content?.attachments?.length > 0 ||
    post?.content?.links?.length > 0 

$: im = post.type == `m.room.message`

$: review = post.content["com.hummingbard.review"] != null &&
    post.content["com.hummingbard.review"] != undefined

$: rating = post.content["com.hummingbard.review"]?.rating
$: reviewTitle = post.content["com.hummingbard.review"]?.title
$: reviewContent = post.content["com.hummingbard.review"]?.content

let loadingReplies = false;

function loadMoreReplies() {
    loadingReplies = true
  fetchMoreReplies().then((res) => {
    console.log(res)
      if(res?.replies && res?.replies.length > 0) {
          res?.replies.forEach(reply => {
              let ind = post.replies.findIndex(x => x.event_id == reply.event_id)
              if(ind  == -1) {
                post.replies = [...post.replies, reply]
              }
          })
      }
      if(res?.unsorted  && res?.unsorted.length > 0) {
          addPosts(res?.unsorted)
      }
  }).then(() => {
    loadingReplies = false
  })
}

async function fetchMoreReplies() {
    let endpoint = `/post/replies/fetch`

    let data = {
        room_id: window.timeline?.room_id,
        event_id: post?.event_id,
        thread_room_id: post?.room_id,
        /*
        room_id: window.timeline?.room_id,
        event_id: window.timeline?.permalinkedPost?.event_id,
        thread_room_id: post?.room_id,
        */
    }

  console.log(data)

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: JSON.stringify(data),
    headers:{
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

let Reply;
let replyLoaded;
function loadReply() {
    if(Reply && !replyLoaded) {
        replyLoaded = true
        return
    }
  import('../new-post/view/view.svelte').then(res => {
    Reply = res.default
    replyLoaded = true
  })
}

function replyToPost() {
    loadReply()
}


function added(e) {
    console.log(e.detail, post.event_id)
    post.replies = [...post.replies, e.detail.post]
    if(reply) {
        dispatch('added', e.detail)
    }
}


function killed() {
    replyLoaded = false
}


function sharePost() {
    let sh = reply || nestedReply || false
    let replyPermalink;
    if(sh) {
        replyPermalink =
            `${window.timeline?.permalinkedPost?.content?.room_path}/${window.timeline?.permalinkedPost?.event_id}/${post.event_id}`
    }
    window.sharePost(post.event_id, (reply || nestedReply || false), replyPermalink)
}

let view;
$: if(view && post.just_added && ( reply || nestedReply )) {
    view.scrollIntoView({ behavior: 'smooth', block: 'center'});
}

$: children = post.unsigned?.children?.['m.reference'] || 0

$: loadMore = post?.replies?.length < children



$: nsfw = (post.content?.nsfw && !identity.preferences.show_nsfw_posts) || false

let nsfwState = true;
function nsfwOff() {
    nsfwState = false
}

$: pi = window.location.pathname.split("/")
$: lastPi = pi[pi.length - 1]

$: replyPermalinked = lastPi === post.event_id

$: showReact = window.timeline?.userFeed || window.timeline?.member
</script>

<div class="po-co fl-co-o" 
class:relative={nsfw}
bind:this={view}
id={post.short_id}
     class:brd-btm={depth == 0}
     class:just-added={post.just_added && (reply || nestedReply)}>


<div class="pi fl-co-o lh-copy brd-b-0" 
     class:brd={embed || nested}
    class:permalinked={replyPermalinked}>


    <div class="flex">
      <div class="mr3">
          <a href={senderHREF} title={post.sender}>
            {#if post.author?.avatar_url?.length > 0}
                <div class={nested ? 'thumbnail-s':'thumbnail'}>
                  <img src="{post.author.avatar_url}" />
                </div>
            {:else}
                {#if nested || reply || nestedReply}
                    <svg class="gr-center" height="22" width="22">
                       <circle cx="11" cy="11" r="11" stroke-width="0" fill="black" />
                    </svg>
                {:else}
                    <svg class="gr-center" height="30" width="30">
                       <circle cx="15" cy="15" r="15" stroke-width="0" fill="black" />
                    </svg>
                {/if}
            {/if}
         </a>
      </div>
      <div class="flex flex-column flex-one gr-default">

        <div class="small gr-start-center">
          <a href={senderHREF}>
            <span class=""><strong>{user}</strong></span>
          </a>

          {#if post.is_article && !nested && !embed}
              posted an article
          {/if}

          {#if review && !nested && !embed}
              posted a review
          {/if}


          {#if !nestedReply && !reply && (!post.content?.share_reply)}
              {#if !window.timeline?.permalink}
                  {#if post?.content?.room_path?.length > 0 && showPath}
                      in <a href={post.content?.room_path}><span class="primary">{post.content.room_path}</span></a>
                  {/if}
              {/if}
          {/if}

          {#if nestedReply || reply}
              <span class="ml3" title={post.date}}>{post.when}</span>
            {/if}

        </div>

        {#if !nestedReply && !reply}
        <div class="small o-90">
          <span title={post.date}}>{post.when}</span>
        </div>
        {/if}

      </div>

      {#if !embed}
      <div class="post-tools flex">
          {#if !nested && !im}
          <div class="perma-link ml3"
              class:gr-center={reply || nestedReply}>
              <a class="" href={permalink}>
                  <svg class="p-t-f pointer" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M10.604 1h4.146a.25.25 0 01.25.25v4.146a.25.25 0 01-.427.177L13.03 4.03 9.28 7.78a.75.75 0 01-1.06-1.06l3.75-3.75-1.543-1.543A.25.25 0 0110.604 1zM3.75 2A1.75 1.75 0 002 3.75v8.5c0 .966.784 1.75 1.75 1.75h8.5A1.75 1.75 0 0014 12.25v-3.5a.75.75 0 00-1.5 0v3.5a.25.25 0 01-.25.25h-8.5a.25.25 0 01-.25-.25v-8.5a.25.25 0 01.25-.25h3.5a.75.75 0 000-1.5h-3.5z"></path></svg>
              </a>
          </div>
          {/if}
          {#if authenticated && !nested}
              <div class="post-menu ml3" 
              class:gr-center={reply || nestedReply}>
                  <Menu id={post.event_id}/>
              </div>
          {/if}
      </div>
    {/if}


    </div>




    <div class=""
         class:post-container={!nestedReply}
         class:post-container-a={nested || reply}>

        {#if post.content.msgtype == 'm.text'}

            <Content post={post} reply={reply || nestedReply} />


            {#if post.shared_post}
              <svelte:self 
                  post={post.shared_post}
                  nested={true}
              />
            {/if}



            {#if post.content?.images?.length > 0}
                <div class={(reply || nestedReply || embed) ? 'flex flex-wrap' : 'post-images'}
                    class:mb3={!post.content?.links?.length || !post.content?.attachments?.length}>
                    {#each post.content.images as image (image.mxc)}
                        <Image image={image} reply={reply || nestedReply || embed} />
                    {/each}
                </div>
            {/if}


            {#if post.content?.links?.length > 0}
            <div class="link-items" class:mb3={!post.content?.attachments?.length}>
                {#each post.content.links as link (link.href)}
                    <Link link={link} reply={reply || nestedReply || embed}/>
                {/each}
            </div>
            {/if}


            {#if post.content?.attachments?.length > 0}
            <div class="attachment-items flex flex-column mb3 pa3">
                {#each post.content.attachments as attachment (attachment.mxc)}
                    <Attachment 
                        attachment={attachment} 
                        single={post.content?.attachments?.length == 1}/>
                {/each}
            </div>
            {/if}

            {#if review}
            <div class="link-item pa3 flex flex-column fs-09 mb3">

                <div class="">
                  <strong>{reviewTitle}</strong>
                </div>
                <div class="mt3">
                    {reviewContent}
                </div>

                <div class="mt3 flex stars">
                    {#each Array(rating) as _, i}
                        <div class="mr1">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25z"></path></svg>
                        </div>
                    {/each}
                    {#each Array(7 - rating) as _, i}
                        <div class="mr1">
                          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25zm0 2.445L6.615 5.5a.75.75 0 01-.564.41l-3.097.45 2.24 2.184a.75.75 0 01.216.664l-.528 3.084 2.769-1.456a.75.75 0 01.698 0l2.77 1.456-.53-3.084a.75.75 0 01.216-.664l2.24-2.183-3.096-.45a.75.75 0 01-.564-.41L8 2.694v.001z"></path></svg>
                        </div>
                    {/each}
                </div>
            </div>
            {/if}

        {/if}

        {#if !reply && !nested && !nestedReply && !im}
          <div class="flex" style="min-height: 26px;" class:mt3={mediaExists}>
            <div class="mr3 gr-center">
                  <a class="" href="{permalink}">
                    {#if post.total_replies > 0}
                      <span class="small o-80 hov-un">{post.total_replies} {post.total_replies > 1 ? 'Replies' : 'Reply'}</span>
                    {:else}
                      <span class="small o-80 hov-un">Reply</span>
                    {/if}
                  </a>
              </div>
              <Reactions id={post.event_id}/>
          </div>
        {/if}

    </div>

    <div class="flex">
        {#if authenticated && (window.timeline?.member || window.timeline?.owner)}
            {#if reply || nestedReply}
                <div class="mr3" class:post-container-a={depth ==0}>
                <span class="small o-70 hov-op pointer" on:click={replyToPost}>Reply</span>
            </div>
            {/if}
        {/if}
        {#if loadMore && depth != 0}
            <div class="gr-center">
                <span class="small hov-un pointer primary"
                      on:click={loadMoreReplies}>{loadingReplies ? 'Loading' : 'Load'} More</span>
            </div>
        {/if}
    </div>


</div>


{#if reply || nestedReply}
    {#if replyLoaded}
        <div class="mh3 brd-tp brd-lr" class:pad-n={reply}>
          <Reply 
            reply={true} 
            eventID={post.event_id}
            on:added={added}
            on:killed={killed}/>
        </div>
    {/if}
{/if}



{#if post.replies?.length > 0 && !embed}
    <div class="flex post-replies">
        {#if depth == 0}
        <div class="pad-n">
        </div>
        {/if}
        <div class="flex-one">
    {#each post.replies as reply (reply.event_id)}
      <svelte:self 
          post={reply}
          nestedReply={true}
          depth={depth + 1}
          expanded={nestedReply ? expanded : expand}
          on:collapse
          on:collapse={collapse}
      />
    {/each}
        </div>
    </div>
{/if}

{#if loadMore && depth == 0}
    <div class="pa2 gr-default">
        <div class="gr-center flex">
            <div class="gr-center">
              <a class="" href={permalink}>
                <span class="small hov-un pointer primary">
                    View More Replies
                </span>
              </a>
            </div>
            <div class="ml2 gr-center">
              <a class="" href={permalink}>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M6.22 3.22a.75.75 0 011.06 0l4.25 4.25a.75.75 0 010 1.06l-4.25 4.25a.75.75 0 01-1.06-1.06L9.94 8 6.22 4.28a.75.75 0 010-1.06z"></path></svg>
              </a>
            </div>
        </div>
    </div>
{/if}

{#if !nested && !embed && !nestedReply && !reply}
{/if}


{#if nsfw && nsfwState}
    <div class="nsfw-mask gr-default brd-btm" class:brd-l={nestedReply}>
        <div class="gr-center small">
            NSFW. <span class="primary pointer" on:click={nsfwOff}>Click to View</span>
        </div>
    </div>
{/if}

</div>



<style>
.just-added {
    outline: 1px solid var(--primary-dark-gray);
}

.n-ic {
    position: absolute;
    top: 50%;
    left: 0;
    right: 0;
}

.permalinked {
    background: var(--primary-lightestest-gray);
}

.brd-b-0 {
    border-bottom: none;
}

.pad-n {
    margin-left: calc(22px + 1rem);
}

@media screen and (max-width: 768px) {
    .pad-n {
        margin-left: 22px;
    }

}

</style>
