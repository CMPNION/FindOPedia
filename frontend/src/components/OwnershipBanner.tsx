import type { Owner } from '../api/articles'

export function OwnershipBanner({ owner }: { owner: Owner | null }) {
  if (!owner) {
    return (
      <div className="alert alert-success">
        Unclaimed — be the first to own this article!
      </div>
    )
  }
  return (
    <div className="alert alert-warning">
      Owned by <strong>{owner.username}</strong> · {new Date(owner.claimed_at).toLocaleDateString()}
    </div>
  )
}
