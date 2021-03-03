<script>
import Link from './link/link.svelte'
import Image from './image/image.svelte'
import Attachment from './attachment/attachment.svelte'
import NSFW from './nsfw/nsfw.svelte'
import Anonymous from './anonymous/anonymous.svelte'
import Article from './article/article.svelte'
import ArticleSettings from './article/settings.svelte'

export let store;
export let allowArticle;

$: isArticle = store?.article?.enabled

$: showArticle = !window.timeline?.is_article

$: profile = window.timeline?.profile

</script>

<div class="flex">

  {#if !isArticle}
    <div class="gr-center">
      <Link store={store} on:addLink/>
    </div>
    <div class="gr-center ml3">
      <Image store={store} on:addImage/>
    </div>
    <div class="gr-center ml3">
      <Attachment store={store} on:addAttachment/>
    </div>
    <div class="gr-center ml3">
      <NSFW store={store} on:toggleNSFW/>
    </div>
  {/if}

  {#if allowArticle}
    <Article 
      store={store} 
      on:toggleArticleOn
      on:toggleArticleOff
      on:killArticleSettings/>
  {/if}

  {#if isArticle}
    <ArticleSettings 
    store={store} 
    on:toggleArticleSettings/>
  {/if}

</div>
