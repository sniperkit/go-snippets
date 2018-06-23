import MovieList from './MovieList'
import { fetchMovies } from '../../Services/Movies'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'

const mapStateToProps = state => ({
    movies: state.movies
})
const mapDispatchToProps = dispatch => bindActionCreators({ fetchMovies }, dispatch)

export default connect(mapStateToProps, mapDispatchToProps)(MovieList)