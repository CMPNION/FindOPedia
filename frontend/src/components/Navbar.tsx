import { NavLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/auth'

export function Navbar() {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <nav className="nav">
      <NavLink to="/" className="nav-brand">FindOPedia</NavLink>
      <div className="nav-spacer" />
      <div className="nav-links">
        {user ? (
          <>
            <NavLink to="/leaderboard" className="nav-link">Leaderboard</NavLink>
            <NavLink to={`/users/${user.username}`} className="nav-link">Collection</NavLink>
            <NavLink to="/settings" className="nav-link">Settings</NavLink>
            <button className="nav-btn-ghost" onClick={handleLogout}>Logout</button>
          </>
        ) : (
          <>
            <NavLink to="/login" className="nav-link">Login</NavLink>
            <NavLink to="/register" className="nav-btn">Register</NavLink>
          </>
        )}
      </div>
    </nav>
  )
}
