<template>
  <div>     
    <Header/>   
    <Container/>
    <Footer/>
  </div>
</template>



<script>
import Header from './components/header.vue'
import Container from './components/container.vue'
import Footer from './components/footer.vue'

export default {
  components:{
    Header,Container,Footer
  }  
}
</script>
