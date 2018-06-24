<template>
    <div> 
        <h3 class="title">More</h3>
        <div class="more">            
            <span>mino built by vuejs,</span>
            <a target="_blank" href="https://github.com/Lemonreds/chrome-extension-mino">@see more.</a>
        </div>
    </div>
</template>


<script>
export default {
   
}
</script>

<style lang="scss">
.more{     
    font-size: 13px;
    color: rgba(0,0,0,.74); 
    line-height: 18px;
    &>a{
        text-decoration: underline;
        color: rgba(0,0,0,1)
    }    
}
</style>

