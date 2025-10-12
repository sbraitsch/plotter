Algorithm Outline:

- Assign all people with plot priority their highest uncontested plot (highest prio for that plot)
- Go through all contested plots
    - RNG a winner
    - Assign prio status to losers
- Continue until all plot priority players are assigned or out of scored plots
- Assign unscored (by plot prio) plots to neighbor prio players
    - Prefer plots they scored
    - Same logic as above
    - How to handle no neighboring plots being left?
- RNG leftover plots among leftover players

Rules:

- Neighbor picks must be bidirectional
- Plot Prio -> Neighbor wish ignored
