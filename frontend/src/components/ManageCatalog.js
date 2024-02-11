import React, { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";

const ManageCatalog = () => {
    const [movies, setMovies] = useState([])
    const { jwtToken } = useOutletContext()
    const { toggleRefresh } = useOutletContext()
    const navigate = useNavigate()

    useEffect(() => {
        toggleRefresh()
    }, [toggleRefresh])

    useEffect(() => {
        if (jwtToken === "") {
            navigate("/login")
            return
        }

        const headers = new Headers()
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", "Bearer " + jwtToken)

        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`/admin/movies`, requestOptions)
            .then(response => response.json())
            .then(data => setMovies(data))
            .catch(error => console.log(error))
    }, [jwtToken])

    return (
        <div>
            <h2>Manage Catalog</h2>
            <hr />

            <table className="table table-striped table-hover">
                <thead>
                    <tr>
                        <th>Movie</th>
                        <th>Release Date</th>
                        <th>Rating</th>
                    </tr>
                </thead>

                <tbody>
                    {movies.map((m) => (
                        <tr key={m.id}>
                            <td>
                                <Link to={`/admin/movies/${m.id}`}>
                                    {m.title}
                                </Link>
                            </td>
                            <td>{m.release_date}</td>
                            <td>{m.mpaa_rating}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

export default ManageCatalog