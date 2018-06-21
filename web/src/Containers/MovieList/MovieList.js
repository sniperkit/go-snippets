import React, { Component } from 'react'
import PropTypes from 'prop-types'
import MovieGrid from '../../Components/MovieGrid'

export default class MovieList extends Component {
    componentWillMount() {
        const { fetchMovies } = this.props

        fetchMovies()
    }

    render() {
        const { loading, movies } = this.props.movies

        return (
            <div>
                {loading && "LOADING"}
                <MovieGrid movies={movies} />
            </div>
        )
    }
}

MovieList.propTypes = {
    movies: PropTypes.object
}

MovieList.defaultProps = {
    movies: {
        loading: false,
        movies: []
    }
}
