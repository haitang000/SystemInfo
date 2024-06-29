document.addEventListener('DOMContentLoaded', function () {
    async function fetchData() {
        try {
            const response = await fetch('/info');
            const data = await response.json();

            // Update status
            document.getElementById('totalMemory').textContent = data.totalMemory;
            document.getElementById('freeMemory').textContent = data.freeMemory;
            document.getElementById('usedMemory').textContent = data.usedMemory;
            document.getElementById('memoryUsage').textContent = data.memoryUsage.toFixed(2);
            document.getElementById('cpuUsage').textContent = data.cpuUsage.toFixed(2);
        } catch (error) {
            console.error('Error fetching system info:', error);
        }
    }

    setInterval(fetchData, 60000); // 每分钟刷新一次数据
    fetchData(); // 初始获取一次数据
});