export default function Cell({ value, falling }) {
  let display = ""

  if (value === 88) display = "X"
  if (value === 79) display = "O"

  return (
    <div className={`cell ${display ? "filled" : ""} ${falling ? "drop" : ""}`}>
      {display}
    </div>
  )
}
