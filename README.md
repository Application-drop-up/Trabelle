# ✈️ trabelle
> Stop jumping between Google Maps, notes, and messaging apps. One-stop travel planning made simple.

---

## 🌟 Overview
`trabelle` is a one-stop web application that allows you to manage all your travel spots, to-dos, and notes on a single integrated map, eliminating the need to constantly switch between multiple tools.

## 👥 Problem Statement (Personas)
This app is designed to solve specific pain points for three different types of travelers:

- **The Planner:** Finds entering tourist spot names tedious and wants auto-complete. Also wants to reuse spot information already registered by other users.
- **The Minimalist:** Uses LINE notes and Google Maps but wants to avoid over-complicated apps with useless features. Needs to keep pre-travel logistics (entry applications, required cards) in one place.
- **The Researcher:** Researches on Instagram/web, saves notes, and drops pins on Google Maps. Finds switching between maps and notes annoying, and finds sharing via messaging apps inconvenient due to the lack of update notifications.

## 🛠️ Tech Stack & Reasoning

### Frontend: React / TypeScript
* **Component-Based Architecture & High Reusability:** Essential for managing the complex, interactive UI required for a map-centric application. UI elements like custom spot cards, color-coded pins, and interactive note panels are built as isolated components and seamlessly reused across different views, significantly improving maintainability.
* **Declarative UI with Virtual DOM:** Instead of imperatively manipulating the DOM (which is highly complex in map applications), React allows us to declare the UI state. The Virtual DOM ensures that only the modified parts of the map or notes are re-rendered, providing a highly performant and smooth user experience.
* **Predictable One-Way Data Flow:** By enforcing a strict top-down data flow (Props), the state of the application remains predictable and easy to debug. This is crucial for tracing how spot data moves from the map interface to the detail panels, preventing hard-to-find bugs as the app grows.
* **Rich Ecosystem:** Leverages mature libraries like `@react-google-maps/api` to integrate efficiently with Google Maps APIs, allowing focus on core product features rather than low-level map implementations.

### Backend: Go
* **High Performance & Concurrency:** Ideal for a scalable API server that handles numerous concurrent requests, particularly for external API calls (like Google Places API) and future real-time collaborative features.
* **Strict Type System:** Ensures type safety from backend to frontend, reducing runtime errors and improving codebase maintainability in the long term.
* **Deployment Efficiency:** Compiles to a single binary for fast and straightforward deployment, making it perfect for a modern, containerized infrastructure.

### Infrastructure & APIs
* **Google Maps API / Google Places API:** Chosen for its industry-standard accuracy and comprehensive global point-of-interest (POI) database, which is critical for providing a seamless auto-complete and mapping experience.

## 📅 Roadmap

### MVP
- [ ] Spot auto-complete (Google Places API)
- [ ] Categorized pin management (Color-coded for breakfast, dinner, etc.)
- [ ] Notes linked to specific spots
- [ ] Read-only URL sharing feature

### v2
- [ ] User authentication & registration
- [ ] Pre-travel essential information (Country-specific guides)
- [ ] Reusing spot information from other users
- [ ] Collaborative plan editing

### v3
- [ ] Real-time update notifications
- [ ] Real-time collaborative editing

---

## 🚫 Out of Scope
To maintain the app's simplicity and focus, the following features are currently out of scope:
- Photo uploading features
- Social media-like posting features
