<template> 
    <div class="item application">   
        <i 
            @click="applicationHandler()"
            class="material-icons"
            >{{option.icon}}</i>
        <transition name="top">
            <ul 
                v-if="visible"
                class="list dialog">
                <li                 
                    v-for="item in application"
                    :key="item.name"
                    >
                    <a :href="item.link">
                        <img :src="item.icon">
                        <span>{{item.name}}</span>
                    </a>                
                </li>  
            </ul> 
        </transition> 
    </div> 
</template>

<script>
import application from '../../api/application.js'

export default {
    data(){
        return{
            application: {}
        }
    },
    props:{
        option:{
            type: Object,
            require: true
        },
        visible:{
            type: Boolean,
            require: true
        }
    },
    created(){
        this.application = application
    },
    methods:{
        close(){
            this.$emit('close','applicationHandler')
        },
        applicationHandler(){
            this.$emit('applicationHandler')
        }
    }
}
</script>

<style lang="scss" scoped>
.application{    
    &>.list{
        width: 320px;
        padding: 28px;    
        & > li{
            width: clac((320px-56px)/3);
            height: 100px;
            display: inline-block;
           
            border: 1px solid transparent;
            border-radius: 2px;
            & > a{ 
                display: block;
                width: 84px;
                height: 84px;
                text-align: center;
                margin: 8px 0;                             
                &> img{
                    width: 64px;
                    height: 64px;
                    padding: 10px;                    
                }
                &> span{
                    text-align: center;                   
                    font-size: 13px;
                    width: 84px;
                    height: 20px;
                    color: #252b32;
                }
                &:hover{
                    opacity: .84;
                    &>span{
                        color: #000;
                    }
                }
            }           
        }
    }   
}
</style>
