import { FETCHING_MOVIES, FETCHING_MOVIES_COMPLETE } from './types'

const initialState = {};

export default function reducer(state = initialState, action) {
    switch (action.type) {
        case FETCHING_MOVIES:
            return { ...state, loading: true, movies: [] }

        case FETCHING_MOVIES_COMPLETE:
            return { ...state, loading: false, movies: action.movies }

        default:
            return state
    }
}
