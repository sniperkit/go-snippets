<template> 
    <div class="search"
        :class = "(commonsites && commonsites.status) ? styles.hasCommonsites : styles.noCommonsites "
    >      
        <div class="icon"
                @click="dialogVisible">
            <img 
                :src="current.icon"
                :alt="current.name">
        </div>      
        <div class="bar">   
            <form :action="current.action">       
                <input type="text" 
                :name="current.inputName" 
                autofocus autocomplete="off" placeholder="Go anywhere.">   
            </form>           
        </div>  

        <Engines v-if="enginesDialog"
                :engines="engines"
                @closeEngines="dialogVisible"
                @changeEngine="changeEngine"
        />

    </div>
</template>

<script>

import * as engineAPI from '../../api/engines'
import Engines from './engines.vue'

import {mapGetters} from 'vuex'

export default {  
    data(){
        return{
            engines : null,
            current : {},
            enginesDialog: false,
            styles:{
                hasCommonsites : '',
                noCommonsites: "noCommonsites"
            }
        }
    },
    created(){
        //get data from api
        this.engines = engineAPI.allEngines
        engineAPI.getEngineIndx.then( idx =>{
            this.current = this.engines[idx]
        })
    },   
    computed:{
         ...mapGetters([     
            'commonsites'
        ])       
    },  
    components:{
        Engines
    },
    methods: {
        // 显示/隐藏 搜索引擎对话框
        dialogVisible() {
            this.enginesDialog = !this.enginesDialog
         },
         // 更改默认的搜素引擎
        changeEngine(idx){
            this.current = this.engines[idx]
            this.dialogVisible()
            engineAPI.setEngineIndx(idx)
        }
    }
};
</script>

<style lang="scss" scoped>
.search {
    width: 640px;
    height: 70px;
    padding: 0 2px;

    display: flex;
    align-items: center;
    
    position: absolute;
    top: 44%;
    left: 50%;
    transform: translate(-50%,-50%);
    z-index: inherit;
    & > .icon{
        width: 55px;
        height: 55px;
        
        padding: 2px;
        background: #fff;
        opacity: .84;

        border: 1px solid #eee;
        border-radius: 4px;

        display: flex;
        align-items: center;
        justify-content: center;

        cursor: pointer;        
        & > img{
            width: 36px;
            height: 36px;
        }
    }
    & > .bar{
        width: 520px;
        height: 55px;
        margin-left: 18px;

        display: inline-block;
        & >form{
            width: 100%;
            height: 100%;
        }
    }  
}
.noCommonsites{
    top: 48%;
    left: 50%;
}


input[type="text"]{
    display: block;
    width: 100%;
    height: 55px;
    padding: 8px 12px;
    font-family: inherit;
    font-size: 16px;
    font-weight: inherit;
    background: rgba(255,255,255,.84);
    outline: none;
    border: 1px solid #eee;    
    border-radius: 2px;
    color: rgba(0,0,0,.74)   
}

</style>
