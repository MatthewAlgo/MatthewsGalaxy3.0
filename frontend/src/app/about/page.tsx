import SubscribeForm from '@/components/SubscribeForm';

export const metadata = {
    title: "About Matei | Matthew's Galaxy",
    description: "Learn about Matei-Alexandru Dinu - Software Engineer with experience at Google, Snowflake, and NXP Semiconductors.",
};

export default function AboutPage() {
    return (
        <div style={{ flex: 1 }}>
            {/* Hero */}
            <section className="about-hero">
                <div className="container">
                    <div className="about-hero-content">
                        <div className="about-avatar">👨‍💻</div>
                        <h1 className="about-name">Matei-Alexandru Dinu</h1>
                        <p className="about-tagline">
                            Software Engineer • Cloud Enthusiast • Tech Explorer
                        </p>
                        <p className="about-location">📍 Warsaw, Poland</p>

                        <div className="about-links">
                            <a href="https://www.linkedin.com/in/username" target="_blank" rel="noopener noreferrer" className="about-social-link">
                                LinkedIn
                            </a>
                            <a href="https://github.com/username" target="_blank" rel="noopener noreferrer" className="about-social-link">
                                GitHub
                            </a>
                            <a href="mailto:admin@example.com" className="about-social-link">
                                Email
                            </a>
                        </div>
                    </div>
                </div>
            </section>

            {/* Bio */}
            <section className="about-section container">
                <div className="about-bio">
                    <h2>About Me</h2>
                    <p>
                        Driven and results-oriented Software Engineer with a proven track record of delivering
                        high-impact projects at top-tier companies, including <strong>Snowflake</strong>, <strong>Google</strong>,
                        and <strong>NXP Semiconductors</strong>.
                    </p>
                    <p>
                        I graduated Computer Science at University of Bucharest with honors, specializing in
                        cloud-native technologies, backend systems, and open-source development. My thesis on
                        &quot;Scalable Cloud Platforms for Generative AI Models&quot; received a perfect 10.0/10 grade.
                    </p>
                    <p>
                        I thrive in fast-paced environments where I can learn, collaborate, and make a tangible impact.
                        Whether it&apos;s building internal tools at Google, automating Kubernetes clusters at Snowflake,
                        or contributing to Matter SDK at NXP, I&apos;m passionate about solving complex problems with elegant solutions.
                    </p>
                </div>
            </section>

            {/* Experience Timeline */}
            <section className="about-experience-section">
                <div className="container">
                    <h2 className="about-section-title">💼 Experience</h2>
                    <div className="about-timeline">
                        <div className="about-timeline-item">
                            <div className="about-timeline-dot"></div>
                            <div className="about-timeline-content">
                                <div className="about-timeline-header">
                                    <h3>Software Engineer Intern</h3>
                                    <span className="about-company">❄️ Snowflake</span>
                                </div>
                                <span className="about-period">Sep 2025 - Dec 2025 • Warsaw, Poland</span>
                                <p>
                                    Part of Cloud Engineering, Container Platform team. Created automation modules in Go
                                    for Kubernetes clusters, facilitating easy developer experience. Wrote unit tests
                                    with 90+% code coverage and created team presentations for new features.
                                </p>
                                <div className="about-tags">
                                    <span>Go</span><span>Kubernetes</span><span>Docker</span><span>Infrastructure</span>
                                </div>
                            </div>
                        </div>

                        <div className="about-timeline-item">
                            <div className="about-timeline-dot"></div>
                            <div className="about-timeline-content">
                                <div className="about-timeline-header">
                                    <h3>Software Engineer Intern</h3>
                                    <span className="about-company">🔷 Google</span>
                                </div>
                                <span className="about-period">Jun 2025 - Sep 2025 • Warsaw, Poland</span>
                                <p>
                                    Part of Google Cloud, Experiments team. Created and released an internal tool that
                                    automates a critical process used by hundreds of clients. Took the project from
                                    concept to official release with real-world impact.
                                </p>
                                <div className="about-tags">
                                    <span>Go</span><span>Cloud</span><span>Internal Tools</span>
                                </div>
                            </div>
                        </div>

                        <div className="about-timeline-item">
                            <div className="about-timeline-dot"></div>
                            <div className="about-timeline-content">
                                <div className="about-timeline-header">
                                    <h3>Software Engineer</h3>
                                    <span className="about-company">🟢 NXP Semiconductors</span>
                                </div>
                                <span className="about-period">Mar 2023 - May 2025 • Bucharest, Romania</span>
                                <p>
                                    Developed full-stack application for Matter SDK in C++ and Angular, officially
                                    integrated into the open-source Matter layer. Created dynamic packet routing
                                    systems, OTA update features, and represented the team at Innovation World Tour 2024.
                                </p>
                                <div className="about-tags">
                                    <span>C++</span><span>Angular</span><span>Python</span><span>IoT</span><span>Matter</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Education */}
            <section className="about-section container">
                <h2 className="about-section-title">🎓 Education</h2>
                <div className="about-education">
                    <div className="about-edu-card">
                        <div className="about-edu-icon">🏛️</div>
                        <div className="about-edu-content">
                            <h3>University of Bucharest</h3>
                            <p className="about-degree">Bachelor of Science - Computer Science</p>
                            <p className="about-edu-period">Aug 2022 - Oct 2025</p>
                            <ul className="about-edu-highlights">
                                <li>🏆 Final Project Grade: <strong>10.0 / 10</strong></li>
                                <li>📊 Last Year Average: <strong>9.83 / 10</strong></li>
                                <li>📚 General Average: <strong>9.25 / 10</strong></li>
                                <li>🎖️ Admitted with honors, full scholarship</li>
                                <li>📝 Thesis: &quot;Scalable Cloud Platform for GenAI Models&quot;</li>
                            </ul>
                        </div>
                    </div>
                </div>
            </section>

            {/* Skills */}
            <section className="about-skills-section">
                <div className="container">
                    <h2 className="about-section-title">🛠️ Skills & Technologies</h2>
                    <div className="about-skills-grid">
                        <div className="about-skill-category">
                            <h4>Languages</h4>
                            <div className="about-skill-list">
                                <span>Go</span><span>Python</span><span>TypeScript</span>
                                <span>C++</span><span>Java</span><span>SQL</span>
                            </div>
                        </div>
                        <div className="about-skill-category">
                            <h4>Cloud & DevOps</h4>
                            <div className="about-skill-list">
                                <span>Kubernetes</span><span>Docker</span><span>AWS</span>
                                <span>GCP</span><span>Linux</span><span>CI/CD</span>
                            </div>
                        </div>
                        <div className="about-skill-category">
                            <h4>Frontend</h4>
                            <div className="about-skill-list">
                                <span>React</span><span>Next.js</span><span>Angular</span>
                                <span>TypeScript</span><span>CSS</span>
                            </div>
                        </div>
                        <div className="about-skill-category">
                            <h4>Backend & Data</h4>
                            <div className="about-skill-list">
                                <span>PostgreSQL</span><span>REST APIs</span><span>gRPC</span>
                                <span>Microservices</span><span>Distributed Systems</span>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Achievements */}
            <section className="about-section container">
                <h2 className="about-section-title">🏆 Achievements</h2>
                <div className="about-achievements">
                    <div className="about-achievement-card">
                        <span className="about-achievement-icon">🥇</span>
                        <h4>1st Place Overall</h4>
                        <p>NitroNLP AI Hackathon</p>
                        <span className="about-achievement-date">Mar 2025</span>
                    </div>
                    <div className="about-achievement-card">
                        <span className="about-achievement-icon">📜</span>
                        <h4>Perfect Thesis Score</h4>
                        <p>10.0/10 - Scalable Cloud Platform for GenAI</p>
                        <span className="about-achievement-date">2025</span>
                    </div>
                    <div className="about-achievement-card">
                        <span className="about-achievement-icon">🎓</span>
                        <h4>Full Scholarship</h4>
                        <p>University of Bucharest</p>
                        <span className="about-achievement-date">2022-2025</span>
                    </div>
                    <div className="about-achievement-card">
                        <span className="about-achievement-icon">🥈</span>
                        <h4>2nd Place</h4>
                        <p>Infoeducatie Programming Contest</p>
                        <span className="about-achievement-date">May 2022</span>
                    </div>
                </div>
            </section>

            {/* Newsletter */}
            <section className="about-section container">
                <SubscribeForm />
            </section>
        </div>
    );
}
