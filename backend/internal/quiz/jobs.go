package quiz

// Jobs is the match catalog. Demand and AIRisk reflect 2026 labor research:
// WEF Future of Jobs 2025 (frontline/care/education growth, graphic design
// decline), PwC 2026 AI Jobs Barometer (AI-skill premium, two-track market),
// and The Guardian's July 2026 AI-safe careers analysis (medicine, teaching,
// hospitality, law).
var Jobs = []Job{
	// Technology & AI
	{Title: "AI / Machine Learning Engineer", HollandCode: "IRE", Category: "Technology", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Build and deploy AI systems. The fastest-growing skill set with a 62% wage premium for AI-skilled workers."},
	{Title: "Cybersecurity Specialist", HollandCode: "IRC", Category: "Technology", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Protect systems and data. Among the top three fastest-growing skill areas through 2030."},
	{Title: "Data Scientist", HollandCode: "ICE", Category: "Technology", Demand: Growing, AIRisk: MediumRisk,
		Blurb: "Turn data into decisions. Core analysis is increasingly AI-assisted, so judgment and framing matter most."},
	{Title: "Software Engineer", HollandCode: "IRC", Category: "Technology", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Still in demand, but transforming fast: AI handles routine code, raising the bar toward system design."},
	{Title: "Data Center Technician", HollandCode: "RIC", Category: "Technology", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Keep the physical infrastructure of AI running. Hands-on and hard to automate."},
	{Title: "AI Ethics / Governance Officer", HollandCode: "ISE", Category: "Technology", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Set the rules for how organizations use AI responsibly. A new standard corporate role."},

	// Healthcare & care
	{Title: "Registered Nurse", HollandCode: "SIR", Category: "Healthcare", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Care roles see some of the largest absolute job growth to 2030. Human presence is the job."},
	{Title: "Physical Therapist", HollandCode: "SRI", Category: "Healthcare", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Hands-on recovery work that AI cannot do. Demographic trends keep demand rising."},
	{Title: "Telehealth Coordinator", HollandCode: "SEC", Category: "Healthcare", Demand: Growing, AIRisk: MediumRisk,
		Blurb: "Manage AI-augmented remote care - one of the fastest-growing healthcare roles."},
	{Title: "Health AI Integration Specialist", HollandCode: "ISR", Category: "Healthcare", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Bridge clinicians and AI tools in real clinical settings. Scarce, hybrid expertise."},
	{Title: "Bioinformatics Scientist", HollandCode: "IRA", Category: "Healthcare", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Analyze genomic data with AI tooling. Deep science plus AI literacy - the premium combination."},
	{Title: "Counselor / Therapist", HollandCode: "SIA", Category: "Healthcare", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Human judgment and trust are the product. Demand for mental-health support keeps climbing."},

	// Education
	{Title: "Teacher (Secondary)", HollandCode: "SAE", Category: "Education", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Education roles grow strongly to 2030, and teaching ranks among the most AI-safe careers."},
	{Title: "Learning Experience Designer", HollandCode: "ASI", Category: "Education", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Design how people learn, often with AI in the loop. Named a top growing role for 2026."},

	// Skilled trades & field work
	{Title: "Electrician", HollandCode: "RCI", Category: "Skilled Trades", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Physical, licensed, and in shortage. Frontline trades see the largest absolute growth to 2030."},
	{Title: "Renewable Energy Technician", HollandCode: "RIC", Category: "Skilled Trades", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Install and maintain solar and wind systems as the energy transition accelerates."},
	{Title: "Plumber / HVAC Technician", HollandCode: "RCE", Category: "Skilled Trades", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Skilled hands-on work no model can replace, with steady demographic demand."},
	{Title: "Construction Manager", HollandCode: "REC", Category: "Skilled Trades", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Coordinate crews and sites. Construction is a top absolute-growth sector to 2030."},
	{Title: "Agritech / Precision Farming Specialist", HollandCode: "RIE", Category: "Skilled Trades", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Apply sensors, drones, and data to food production - a rising sustainability sector."},

	// Business & consulting
	{Title: "Fractional Consultant (CFO/CMO/AI)", HollandCode: "EIS", Category: "Business", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Sell your expertise across several companies. The fastest-growing pivot destination in 2026."},
	{Title: "Product Manager", HollandCode: "EIS", Category: "Business", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Own what gets built and why. Judgment-heavy, but AI is compressing the routine parts."},
	{Title: "Entrepreneur / Founder", HollandCode: "EAS", Category: "Business", Demand: Stable, AIRisk: LowRisk,
		Blurb: "AI lowers the cost of building, which raises the value of people who ship businesses."},
	{Title: "Sales Manager", HollandCode: "ECS", Category: "Business", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Relationship-driven selling stays human; prospecting and admin are being automated."},
	{Title: "Project Manager", HollandCode: "ECS", Category: "Business", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Coordination and accountability across teams. AI eats status reports, not responsibility."},
	{Title: "Operations / Supply Chain Manager", HollandCode: "ECR", Category: "Business", Demand: Growing, AIRisk: MediumRisk,
		Blurb: "E-commerce logistics keeps growing; AI forecasting makes operators more valuable, not less."},
	{Title: "HR Manager", HollandCode: "SEC", Category: "Business", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "People judgment, conflict, and culture stay human even as screening automates."},
	{Title: "Financial Analyst", HollandCode: "CIE", Category: "Business", Demand: Stable, AIRisk: HighRisk,
		Blurb: "Modeling is increasingly automated - analysts who direct AI and own the narrative stay valuable."},
	{Title: "Accountant", HollandCode: "CEI", Category: "Business", Demand: Stable, AIRisk: HighRisk,
		Blurb: "Bookkeeping is automating quickly; advisory and audit judgment are the durable parts."},
	{Title: "Lawyer", HollandCode: "EIS", Category: "Business", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Ranked among AI-resilient careers: accountability and advocacy stay human, research automates."},

	// Creative & media
	{Title: "UX Designer", HollandCode: "AIE", Category: "Creative", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Product thinking and research keep UX valuable while production design automates."},
	{Title: "Graphic Designer", HollandCode: "AER", Category: "Creative", Demand: Declining, AIRisk: HighRisk,
		Blurb: "Honest read: WEF projects decline as generative tools absorb production work. Direction and brand strategy hold value."},
	{Title: "Content Strategist", HollandCode: "AES", Category: "Creative", Demand: Stable, AIRisk: HighRisk,
		Blurb: "Pure writing output is commoditizing; strategy, voice, and editorial judgment are the moat."},

	// Hospitality & service
	{Title: "Hospitality Manager", HollandCode: "ESR", Category: "Hospitality", Demand: Growing, AIRisk: LowRisk,
		Blurb: "Hotels and experiences rank among AI-safe careers - people pay for human hosts."},
	{Title: "Executive / Personal Assistant", HollandCode: "CSE", Category: "Hospitality", Demand: Stable, AIRisk: MediumRisk,
		Blurb: "Varied, trust-based support work; scheduling automates but discretion does not."},
}
