# Airline Crew Scheduling using Tabu Search

This project focuses on optimizing pilot assignments to flights within an airline's operations, using the Tabu Search algorithm. Synthetic data simulates real airline operations, including a predefined set of cities or airports and a set of available pilots.

The goal is to efficiently assign pilots to flights while minimizing operational costs and maintaining a realistic set of constraints. The key objectives include:

Allowing pilots to depart from any airport for their first flight, but restricting subsequent departures to Riga or Liepaja.
Ensuring pilots continue their journey from their current airport, minimizing "transfers" between different airports.
Reducing the number of airplane changes for pilots during their assignments.
The current approach uses penalty-based optimization, penalizing the following violations:
- Pilot departs from an airport different from their current location (-100 points).
- Pilot transfers between airports during their flight sequence (-250 points).
- Pilot changes airplanes (-50 points).
- A bonus is awarded for staying on the same airplane (+10 points).

While the penalty for using a high number of pilots is not yet implemented, the algorithm has shown promising results. It achieves a 30-40% reduction in penalties compared to random assignments, and pilots tend to stick to one airplane whenever possible.

Testing has been performed by generating random scenarios with various seed numbers, keeping the number of pilots, cities, and flights fixed to maintain complexity. Future improvements will focus on implementing a mechanism to discard unnecessary pilots, further reducing operational costs.