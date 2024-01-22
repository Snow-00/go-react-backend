import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

const Movie = () => {
    const [movie, setMovie] = useState({})  // this uses empty object
    let { id } = useParams()

    useEffect(() => {
        let myMovie = {
            id: 1,
            title: "Highlander",
            release_date: "1986-03-07",
            runtime: 116,  // in minutes
            mpaa_rating: "R",
            description: "Some long desc",
        }

        setMovie(myMovie)
    }, [id])  // always render when the id is changed

    return (
        <div>
            <h2>Movie: {movie.title}</h2>
            <small><em>{movie.release_date}, {movie.runtime} minutes, Rated {movie.mpaa_rating}</em></small>
            <hr />
            
            <p>{movie.description}</p>
        </div>
    )
}

export default Movie