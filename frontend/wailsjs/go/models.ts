export namespace domain {
	
	export class AppConfig {
	    steamCmdPath: string;
	    outputDir: string;
	    autoUpdate: boolean;
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.steamCmdPath = source["steamCmdPath"];
	        this.outputDir = source["outputDir"];
	        this.autoUpdate = source["autoUpdate"];
	        this.username = source["username"];
	    }
	}
	export class DownloadJob {
	    id: string;
	    workshopId: string;
	    appId: string;
	    title?: string;
	    previewUrl?: string;
	    status: string;
	    progress: number;
	    errorMsg?: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    startedAt?: any;
	    // Go type: time
	    finishedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new DownloadJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workshopId = source["workshopId"];
	        this.appId = source["appId"];
	        this.title = source["title"];
	        this.previewUrl = source["previewUrl"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.errorMsg = source["errorMsg"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.finishedAt = this.convertValues(source["finishedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

