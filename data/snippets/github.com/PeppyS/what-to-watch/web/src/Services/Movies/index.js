import { FETCHING_MOVIES, FETCHING_MOVIES_COMPLETE } from './types'
import MoviesAPI from '../../API/Movies'
import { wait } from '../../Utils/Timer'

export const fetchMovies = () => async dispatch => {
    dispatch({
        type: FETCHING_MOVIES
    })

    const response = await MoviesAPI.getMovies()

    // Simulate 1 second loading time
    await wait(1000)

    dispatch({
        type: FETCHING_MOVIES_COMPLETE,
        movies: response.data.movies
    })
}
