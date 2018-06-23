import React, { Component } from 'react';
import MovieList from './Containers/MovieList'
import logo from './Images/logo.svg';
import material from './material'
import './Style/App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">what to watch</h1>
        </header>
        <MovieList />
      </div>
    );
  }
}

export default material(App)
