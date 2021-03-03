<script>
export let room;

function alias(alias) {
    let x = alias.substring(1)
    x = x.replace(`:${window.location.hostname}`, "")
    x = x.replaceAll("_", "/")
    return x
}

</script>

<a href="/{alias(room.room_alias)}">
<div class="pv1 flex room-item ph1">
    {#if room.avatar?.length > 0}
        <div class="thumbnail-s mr3">
          <img src="{room.avatar}" />
        </div>
    {:else}
        <div class="mr3">
            <svg class="gr-center" height="22" width="22">
               <circle cx="11" cy="11" r="11" stroke-width="0" fill="black" />
            </svg>
        </div>
    {/if}
    <div class="gr-center flex-one ov-x-h" title={room.room_alias}>
            <div class="primary fs-09 ov-x-h t-ov-h mr2">{alias(room.room_alias)}</div>
    </div>
</div>
</a>

<style>
.room-item {
    transition: 0.1s;
    border-radius: 500px;
}
.room-item:hover {
    background: var(--primary-grayish);
}
</style>
