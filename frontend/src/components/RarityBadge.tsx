export function RarityBadge({ tier }: { tier: string }) {
  return <span className={`badge badge-${tier}`}>{tier}</span>
}
