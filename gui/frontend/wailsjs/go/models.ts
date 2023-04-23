export namespace appvuln {
	
	export class VulnItem {
	    type: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new VulnItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.path = source["path"];
	    }
	}
	export class Info {
	    path: string;
	    executable_path: string;
	    injectable: boolean;
	    dylibs: VulnItem[];
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.executable_path = source["executable_path"];
	        this.injectable = source["injectable"];
	        this.dylibs = this.convertValues(source["dylibs"], VulnItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

export namespace macho {
	
	export class CodeSign {
	    id: string;
	    team_id: string;
	    flags: number;
	    flags_string: string;
	    runtime_version: string;
	    entitlements: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new CodeSign(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.team_id = source["team_id"];
	        this.flags = source["flags"];
	        this.flags_string = source["flags_string"];
	        this.runtime_version = source["runtime_version"];
	        this.entitlements = source["entitlements"];
	    }
	}
	export class Dylib {
	    name: string;
	    time: number;
	    current_version: string;
	    compat_version: string;
	
	    static createFrom(source: any = {}) {
	        return new Dylib(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.time = source["time"];
	        this.current_version = source["current_version"];
	        this.compat_version = source["compat_version"];
	    }
	}
	export class Dylibs {
	    dylinker: string;
	    rpaths: string[];
	    loads: Dylib[];
	    weaks: Dylib[];
	
	    static createFrom(source: any = {}) {
	        return new Dylibs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dylinker = source["dylinker"];
	        this.rpaths = source["rpaths"];
	        this.loads = this.convertValues(source["loads"], Dylib);
	        this.weaks = this.convertValues(source["weaks"], Dylib);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class MachOInfo {
	    magic: string;
	    cpu: string;
	    type: string;
	    dylibs: Dylibs;
	    codesign: CodeSign;
	
	    static createFrom(source: any = {}) {
	        return new MachOInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.magic = source["magic"];
	        this.cpu = source["cpu"];
	        this.type = source["type"];
	        this.dylibs = this.convertValues(source["dylibs"], Dylibs);
	        this.codesign = this.convertValues(source["codesign"], CodeSign);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

