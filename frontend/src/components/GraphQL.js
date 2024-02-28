import React, { useEffect, useState } from "react";
import Input from './form/Input'
import { Link } from "react-router-dom";

const GraphQL = () => {
    // set up stateful vars
    const [movies, setMovies] = useState([])
    const [searchTerm, setSearchTerm] = useState("")
    const [fullList, setFullList] = useState([])

    // perform a search
    const performSearch = () => {

    }

    const handleChange = (event) => {

    }

    // useEffect
    useEffect(() => {
        const payload = `
        {
            list {
                id
                title
                runtime
                release_date
                mpaa_rating
            }
        }`

        const headers = new Headers()
        headers.append("Content-Type", "application/graphql")

        const requestOptions = {
            method: "POST",
            headers: headers,
            body: payload,
        }

        fetch(`/graph`, requestOptions)
            .then(response => response.json())
            .then(response => {
                let theList = Object.values(response.data.list)
                setMovies(theList)
                setFullList(theList)
            })
            .catch(err => console.log(err))
    }, [])

    return (
        <div>
            <h2>GraphQL</h2>
            <hr />

            <form onSubmit={handleChange}>
                <Input
                    title={"Search"}
                    type={"search"}
                    name={"search"}
                    className={form-control}
                    value={searchTerm}
                    onChange={handleChange} />
            </form>

            {movies ? (
                <table className="table table-striped table-hover">
                    <thead>
                        <tr>
                            <th>Movie</th>
                            <th>Release_date</th>
                            <th>Rating</th>
                        </tr>
                    </thead>

                    <tbody>
                        {movies.map((m) => (
                            <tr key={m.id}>
                                <td>
                                    <Link to={`/movies/${m.id}`}>
                                        {m.title}
                                    </Link>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            ) : (
                <p> No movies (yet)</p>
            )}
        </div>
    )
}

export default GraphQL