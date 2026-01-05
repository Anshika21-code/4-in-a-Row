export default function Cell({ value }) {
  return (
    <div className="cell">
      {value === "X" && <div className="disc blue" />}
      {value === "O" && <div className="disc red" />}
    </div>
  );
}
