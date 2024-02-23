import React, { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";

const ManageCatalog = () => {
    const [movies, setMovies] = useState([])
    const { jwtToken } = useOutletContext()
    const { setJwtToken } = useOutletContext()
    const { toggleRefresh } = useOutletContext()
    const navigate = useNavigate()

    const fetchCatalog = useCallback(async () => {
        const headers = new Headers()
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", "Bearer " + jwtToken)

        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        try {
            let response = await fetch(`/admin/movies`, requestOptions)
            if (response.status === 401) {
                toggleRefresh()
                    .then(token => setJwtToken(token))
                    .catch(() => {
                        setJwtToken("")
                        navigate("/login")
                    })
            } else if (response.status === 200) {
                let data = await response.json()
                setMovies(data)
            }
        }
        catch(error) {
            console.log(error)
        }
    }, [jwtToken, toggleRefresh])
    
    useEffect(() => {
        fetchCatalog()
    }, [fetchCatalog])

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
                                <Link to={`/admin/movie/${m.id}`}>
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