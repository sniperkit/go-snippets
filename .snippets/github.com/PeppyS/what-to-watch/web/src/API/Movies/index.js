import API from '../../API'

const getMovies = () => {
    return API.get('/movies')
}

export default {
    getMovies
}
