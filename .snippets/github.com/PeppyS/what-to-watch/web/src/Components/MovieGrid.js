import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import GridListTileBar from '@material-ui/core/GridListTileBar';
import ListSubheader from '@material-ui/core/ListSubheader';
import IconButton from '@material-ui/core/IconButton';
import InfoIcon from '@material-ui/icons/Info';

const styles = theme => ({
    root: {
        display: 'flex',
        flexWrap: 'wrap',
        justifyContent: 'space-around',
        overflow: 'hidden',
        backgroundColor: theme.palette.background.paper,
    },
    gridList: {
        width: 500,
        height: 450,
    },
    icon: {
        color: 'rgba(255, 255, 255, 0.54)',
    },
});

const TitlebarGridList = (props) => {
    const { classes, movies } = props;

    console.log(props)

    return (
        <div className={classes.root}>
            <GridList cellHeight={180} cols={1} spacing={50} className={classes.gridList} style={{ height: 'auto' }}>
                <GridListTile key="Subheader" cols={1} style={{ height: 'auto' }}>
                    <ListSubheader component="div">Trending Movies</ListSubheader>
                </GridListTile>
                {movies.map(movie => (
                    <GridListTile key={movie.title} cols={1}>
                        <img src={movie.image} alt={movie.title} />
                        <GridListTileBar
                            title={movie.title}
                            subtitle={<span>TODO</span>}
                            actionIcon={
                                <IconButton className={classes.icon}>
                                    <InfoIcon />
                                </IconButton>
                            }
                        />
                    </GridListTile>
                ))}
            </GridList>
        </div>
    );
}

TitlebarGridList.propTypes = {
    classes: PropTypes.object.isRequired,
    movies: PropTypes.array
}

TitlebarGridList.defaultProps = {
    movies: []
}

export default withStyles(styles)(TitlebarGridList);
