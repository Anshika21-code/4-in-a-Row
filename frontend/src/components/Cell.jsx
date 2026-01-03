export default function Cell({ value }) {
  let className = "cell"
  if (value === "X") className += " red drop"
  if (value === "O") className += " yellow drop"

  return <div className={className}></div>
}
