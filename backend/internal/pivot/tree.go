package pivot

// RootID is where every thread starts.
const RootID = "driver"

// Nodes is the decision tree. Answer labels double as stored step values,
// so treat them as stable identifiers once shipped.
var Nodes = []Node{
	{
		ID:       "driver",
		Question: "What's really driving the urge to change?",
		Options: []Option{
			{Label: "I'm burnt out - I need more time and energy for life", Next: "hours-fix"},
			{Label: "The work itself no longer fits me", Next: "what-broke"},
			{Label: "I want more autonomy and ownership over my work", Next: "autonomy-kind"},
			{Label: "My field feels at risk (AI, layoffs, decline)", Next: "risk-level"},
		},
	},

	// Burnout / hours branch
	{
		ID:       "hours-fix",
		Question: "If you had 20% more free time, would your current job be fine?",
		Options: []Option{
			{Label: "Yes - the job is fine, it's the hours", Next: "employer-flex"},
			{Label: "No - even less of this job wouldn't fix it", Next: "what-broke"},
		},
	},
	{
		ID:       "employer-flex",
		Question: "How open is your employer to flexible arrangements?",
		Options: []Option{
			{Label: "Open - there are precedents for part-time or 4-day weeks", Outcome: "reduce-hours"},
			{Label: "Rigid - or I haven't dared to ask yet", Next: "trade-income"},
		},
	},
	{
		ID:       "trade-income",
		Question: "Would you trade a meaningful slice of income for time?",
		Options: []Option{
			{Label: "Yes - time matters more right now", Outcome: "portfolio-career"},
			{Label: "Not really - I need the full income", Outcome: "recharge-first"},
		},
	},

	// Work-fit branch
	{
		ID:       "what-broke",
		Question: "What part of the work no longer fits?",
		Options: []Option{
			{Label: "The day-to-day tasks themselves", Next: "role-nearby"},
			{Label: "The company or culture - the craft is still right", Outcome: "internal-move"},
			{Label: "The whole field or direction", Next: "know-field"},
		},
	},
	{
		ID:       "role-nearby",
		Question: "Is there a version of your role nearby that excites you - a different team, specialty, or product?",
		Options: []Option{
			{Label: "Yes - I can name one", Outcome: "internal-move"},
			{Label: "No - I've looked", Next: "know-field"},
		},
	},
	{
		ID:       "know-field",
		Question: "Do you know which field you'd move toward?",
		Options: []Option{
			{Label: "Yes - fairly clearly", Outcome: "switch-field"},
			{Label: "No - that's exactly what I'm trying to figure out", Outcome: "explore-first"},
		},
	},

	// Autonomy branch
	{
		ID:       "autonomy-kind",
		Question: "What kind of autonomy are you after?",
		Options: []Option{
			{Label: "Choosing my clients, projects, and hours", Next: "expertise-demand"},
			{Label: "Building something of my own from scratch", Outcome: "entrepreneurship"},
		},
	},
	{
		ID:       "expertise-demand",
		Question: "Do people already come to you (or pay you) for your expertise?",
		Options: []Option{
			{Label: "Yes - I'm the person others ask for advice", Outcome: "consulting-now"},
			{Label: "Not yet - I'm good but not known for it", Next: "build-proof"},
		},
	},
	{
		ID:       "build-proof",
		Question: "Could you build public proof of your expertise (cases, portfolio, talks) within your current job over the next 6 months?",
		Options: []Option{
			{Label: "Yes - there's room for that", Outcome: "consulting-runway"},
			{Label: "No - my job leaves no room", Outcome: "portfolio-career"},
		},
	},

	// Risk branch
	{
		ID:       "risk-level",
		Question: "How exposed is your role to automation right now?",
		Options: []Option{
			{Label: "Tasks around me are already being automated", Next: "stay-domain"},
			{Label: "It's more a vague worry than a reality yet", Outcome: "upskill-ai"},
		},
	},
	{
		ID:       "stay-domain",
		Question: "Do you want to stay in this field?",
		Options: []Option{
			{Label: "Yes - I like the domain, not the threat", Outcome: "upskill-ai"},
			{Label: "No - I'd rather move somewhere more durable", Next: "know-field"},
		},
	},
}

// Outcomes carry the actual guidance. WhyNow lines cite the 2026 research
// this product is built on (WEF, PwC, SuperCareer, TieTalent, BTG).
var Outcomes = []Outcome{
	{
		ID:      "reduce-hours",
		Path:    "reduce-hours",
		Title:   "Negotiate reduced hours where you are",
		Tagline: "Keep the job, reclaim your time.",
		WhyNow:  "Flexible arrangements are normal enough in 2026 that precedents exist in most organizations - the leverage is asking with a business case, not apologizing.",
		Plan: []string{
			"Audit two weeks of your workload: what only you can do vs what can be dropped or delegated",
			"Draft the proposal as a business case: output you'll protect, coverage plan, review date",
			"Propose a 3-month trial at 80% (or a 4-day week) rather than a permanent change",
			"Anchor on precedents in your company or industry",
			"Agree upfront how success will be measured at the review date",
		},
		Resources: []Resource{
			{Title: "Making the most of a mid-career pivot (Georgia Tech)", URL: "https://pe.gatech.edu/industry-trends/making-the-most-mid-career-pivot"},
			{Title: "How to pivot careers in 2026 (Employment Hero)", URL: "https://employmenthero.com/uk/blog/how-to-change-career-2026/"},
		},
	},
	{
		ID:      "portfolio-career",
		Path:    "reduce-hours",
		Title:   "Build a portfolio career",
		Tagline: "Several income streams, one life that fits.",
		WhyNow:  "What was once labeled freelancing is now a structured, respected workforce model - fractional and multi-client careers are mainstream in 2026 (TieTalent).",
		Plan: []string{
			"List every skill someone has ever paid you (or thanked you) for",
			"Pick one anchor income stream (part-time role or retainer) to cover your baseline",
			"Add one experimental stream - freelance projects, teaching, a productized service",
			"Set a 6-month runway budget so experiments don't create panic",
			"Review quarterly: double down on what compounds, drop what drains",
		},
		Resources: []Resource{
			{Title: "Career change in 2026: why professionals are pivoting (TieTalent)", URL: "https://medium.com/tietalent-com/career-change-in-2026-why-professionals-are-pivoting-and-where-the-real-opportunities-are-97bc293b53d5"},
			{Title: "Career pivot guide 2026 (SuperCareer)", URL: "https://www.supercareer.co/blog/career-pivot-guide-2026"},
		},
	},
	{
		ID:      "recharge-first",
		Path:    "recharge",
		Title:   "Recharge first, decide second",
		Tagline: "Burnout makes every option look wrong.",
		WhyNow:  "Decisions made from depletion optimize for escape, not fit. Stabilize energy first; the pivot question will look different in eight weeks.",
		Plan: []string{
			"Check what leave you actually have: PTO, sabbatical policy, unpaid leave, medical options",
			"Take the longest real break you can - and protect it completely",
			"Fix the one recovery basic that's most broken (sleep, exercise, or boundaries)",
			"Set a decision date 6-8 weeks out; do not decide before it",
			"Revisit this exploration then - fork this thread and answer again with fresh eyes",
		},
		Resources: []Resource{
			{Title: "Making the most of a mid-career pivot (Georgia Tech)", URL: "https://pe.gatech.edu/industry-trends/making-the-most-mid-career-pivot"},
		},
	},
	{
		ID:      "internal-move",
		Path:    "within-field",
		Title:   "Change within your field",
		Tagline: "Same craft, better container.",
		WhyNow:  "A move that keeps your domain expertise compounds it instead of resetting it - and internal or same-field moves are the fastest pivots to execute.",
		Plan: []string{
			"Name precisely what must change: team, manager, product, company, or specialty",
			"Map three concrete target roles where your expertise transfers at full value",
			"Talk to two people already in each target role before applying anywhere",
			"Update your resume for the target role, not your past one - mirror its language",
			"Set a 90-day timeline: conversations in month 1, applications in month 2, decisions in month 3",
		},
		Resources: []Resource{
			{Title: "Best careers to switch to in 2026 (Revarta)", URL: "https://www.revarta.com/blog/best-careers-to-switch-to-2026"},
			{Title: "Build a resume for the target role", URL: "/resumes/new"},
			{Title: "Search openings in your field", URL: "/jobs"},
		},
	},
	{
		ID:      "switch-field",
		Path:    "switch-out",
		Title:   "Switch fields with a structured plan",
		Tagline: "Deliberate beats dramatic.",
		WhyNow:  "82% of career changers who followed a structured plan landed roles within 6 months; successful pivoters spend 3-9 months in deliberate preparation (SuperCareer, 2026).",
		Plan: []string{
			"Write a transferable-skills audit: what you do that the new field also pays for",
			"Close the top skill gap with one focused course or project - not five",
			"Build one public artifact in the new field (project, case study, writing)",
			"Do five informational interviews before your first application",
			"Tailor your resume to the new field's vocabulary and run the keyword check against real postings",
		},
		Resources: []Resource{
			{Title: "Career pivot guide 2026 (SuperCareer)", URL: "https://www.supercareer.co/blog/career-pivot-guide-2026"},
			{Title: "The career pivot playbook (Education Direct)", URL: "https://www.education-direct.com/blog/career-pivot-2026"},
			{Title: "What is the smartest career pivot in 2026? (mba.com)", URL: "https://www.mba.com/business-school-and-careers/career-possibilities/what-is-the-smartest-career-pivot-in-2026"},
			{Title: "Build a field-switch resume", URL: "/resumes/new"},
		},
	},
	{
		ID:      "explore-first",
		Path:    "switch-out",
		Title:   "Explore before you leap",
		Tagline: "Direction first, motion second.",
		WhyNow:  "\"Too old to pivot?\" is the most common anxiety in career forums - the antidote is treating exploration as a project with a deadline, not an identity crisis.",
		Plan: []string{
			"Take the personality test here to get a research-weighted shortlist of fitting fields",
			"Pick the top two results and read a week of real job postings in each",
			"Do two informational interviews per candidate field",
			"Score each field on fit, demand outlook, and entry cost - then commit to one",
			"Fork this thread and answer again once you know your target field",
		},
		Resources: []Resource{
			{Title: "Take the personality test", URL: "/test"},
			{Title: "Future of Jobs Report 2025 (WEF)", URL: "https://www.weforum.org/publications/the-future-of-jobs-report-2025/digest/"},
			{Title: "Which jobs will help you thrive (The Guardian)", URL: "https://www.theguardian.com/money/2026/jul/11/ai-work-jobs-future-medicine-teaching-hotels-law"},
		},
	},
	{
		ID:      "consulting-now",
		Path:    "consultancy",
		Title:   "Go fractional or consulting now",
		Tagline: "Your expertise is the product.",
		WhyNow:  "Fractional executives and consultants are the fastest-growing pivot destination of 2026 - and clients value people who have actually done the work (Business Talent Group).",
		Plan: []string{
			"Define one sharp offer: who you help, with what problem, at what outcome",
			"Package two past wins as one-page case studies - demonstrable outputs beat pedigree",
			"Set your rate: your old daily salary cost x 2-3 is the standard starting band",
			"Land the first client from your existing network before building any brand",
			"Decide your capacity model upfront: how many clients, which days, what's off-limits",
		},
		Resources: []Resource{
			{Title: "Career change to consulting after corporate (Business Talent Group)", URL: "https://resources.businesstalentgroup.com/btg-blog/career-change-consulting-after-corporate"},
			{Title: "The smart way to pivot into consulting (LinkedIn/Jessi Hempel)", URL: "https://www.linkedin.com/pulse/more-freedom-money-smart-way-pivot-consulting-jessi-hempel-ht9oc"},
			{Title: "Consulting for career changers (CaseBasix)", URL: "https://www.casebasix.com/pages/consulting-recruiting-career-changers"},
		},
	},
	{
		ID:      "consulting-runway",
		Path:    "consultancy",
		Title:   "Build your consulting runway",
		Tagline: "Six months of proof, then the leap.",
		WhyNow:  "In 2026, demonstrable outputs - case studies, portfolio projects, visible results - carry real weight with clients and hiring managers alike. Build them before you jump.",
		Plan: []string{
			"Volunteer for the projects in your job that produce measurable, tellable results",
			"Write up one case study per quarter - problem, action, quantified outcome",
			"Start being visible where your future clients look: one post or talk per month",
			"Take one small paid side engagement to test the market (check your contract first)",
			"Set a go/no-go date 6 months out with a concrete trigger: first client or first retainer",
		},
		Resources: []Resource{
			{Title: "Career change to consulting after corporate (Business Talent Group)", URL: "https://resources.businesstalentgroup.com/btg-blog/career-change-consulting-after-corporate"},
			{Title: "Career pivot guide 2026 (SuperCareer)", URL: "https://www.supercareer.co/blog/career-pivot-guide-2026"},
		},
	},
	{
		ID:      "entrepreneurship",
		Path:    "consultancy",
		Title:   "Start your own thing",
		Tagline: "AI lowered the cost of building - judgment is the moat.",
		WhyNow:  "Building has never been cheaper: AI compresses the cost of shipping, which raises the value of people who pick the right thing to ship.",
		Plan: []string{
			"Start with a problem you've personally paid to have solved - or been paid to solve",
			"Talk to ten potential customers before building anything",
			"Ship the smallest sellable version within 60 days",
			"Keep your job until revenue covers your baseline - runway kills more startups than competition",
			"Pick one distribution channel and go deep instead of being everywhere",
		},
		Resources: []Resource{
			{Title: "Career change in 2026 (TieTalent)", URL: "https://medium.com/tietalent-com/career-change-in-2026-why-professionals-are-pivoting-and-where-the-real-opportunities-are-97bc293b53d5"},
		},
	},
	{
		ID:      "upskill-ai",
		Path:    "within-field",
		Title:   "Upskill into the AI track of your field",
		Tagline: "Don't outrun the wave - surf it.",
		WhyNow:  "Workers with AI skills earn a 62% wage premium and AI-skill postings grew 144% year-over-year (PwC 2026). The professionalised track - human judgment plus AI leverage - is growing fastest.",
		Plan: []string{
			"List the three most automatable tasks in your role - those are your delegation targets, not your identity",
			"Learn the AI tools your industry actually uses (not generic courses): 30 minutes a day for 90 days",
			"Become the person who introduces one AI workflow to your team",
			"Add the hybrid keywords to your profile and resume - 'AI-assisted X' roles are where demand is",
			"Re-run the job search quarterly to watch how your field's postings evolve",
		},
		Resources: []Resource{
			{Title: "2026 Global AI Jobs Barometer (PwC)", URL: "https://www.pwc.com/gx/en/services/ai/ai-jobs-barometer.html"},
			{Title: "AI workforce trends 2026 (Gloat)", URL: "https://gloat.com/blog/ai-workforce-trends/"},
			{Title: "Search AI-track roles in your field", URL: "/jobs"},
		},
	},
}

var NodeByID = func() map[string]Node {
	m := make(map[string]Node, len(Nodes))
	for _, n := range Nodes {
		m[n.ID] = n
	}
	return m
}()

var OutcomeByID = func() map[string]Outcome {
	m := make(map[string]Outcome, len(Outcomes))
	for _, o := range Outcomes {
		m[o.ID] = o
	}
	return m
}()
